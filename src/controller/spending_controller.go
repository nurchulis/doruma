package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SpendingController struct {
	SpendingService service.SpendingService
}

func NewSpendingController(spendingService service.SpendingService) *SpendingController {
	return &SpendingController{
		SpendingService: spendingService,
	}
}

func (sc *SpendingController) CreateSpending(c *fiber.Ctx) error {
	webhookURL := "https://n8n-u3amkcsfjd0u.cica.sumopod.my.id/webhook/2f1957b5-18d7-44ae-b66e-123c52e88cff"
	sessionUserID := c.Get("session_user_id")
	authHeader := c.Get("Authorization")

	form, err := c.MultipartForm()
	if err != nil && err != fiber.ErrUnprocessableEntity {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form data"})
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if form != nil && form.File != nil && len(form.File["file"]) > 0 {
		fileHeader := form.File["file"][0]
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot open file"})
		}
		defer file.Close()

		filename := fileHeader.Filename
		ext := filepath.Ext(filename)
		mimeType := fileHeader.Header.Get("Content-Type")

		if mimeType == "" {
			mimeType = mime.TypeByExtension(ext)
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}
		}

		partHeader := make(textproto.MIMEHeader)
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
		partHeader.Set("Content-Type", mimeType)

		part, err := writer.CreatePart(partHeader)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create form file"})
		}

		if _, err := io.Copy(part, file); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot copy file"})
		}
	} else if form != nil && form.Value != nil && len(form.Value["text"]) > 0 {
		text := form.Value["text"][0]
		if err := writer.WriteField("text", text); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot write text field"})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Either text or file must be provided"})
	}

	writer.Close()

	req, err := http.NewRequest("POST", webhookURL, &buf)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create request"})
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if sessionUserID != "" {
		req.Header.Set("session_user_id", sessionUserID)
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send webhook"})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	// Parse response body to get category

	var wr response.WebhookResponse
	_ = json.Unmarshal(body, &wr) // ignore error, wr.Category will be empty if not found

	now := time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")

	createSpending := &validation.CreateSpending{
		UserSessionID: sessionUserID,
		Category:      wr.Category,
		CategoryID:    "95eb84d2-0a32-4aef-b6c2-bfb5bbc686f5",
		Amount:        float64(wr.Total),
		Name:          wr.Used,
		IsConfirm:     true,
		Datetime:      now,
	}
	spending, err := sc.SpendingService.CreateSpending(c, createSpending)
	if err != nil {
		return err
	}
	// Unmarshal body to map
	var respMap map[string]interface{}
	if err := json.Unmarshal(body, &respMap); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse webhook response"})
	}
	// Inject spending ID
	respMap["id"] = spending.ID

	return c.Status(fiber.StatusCreated).
		JSON(response.SuccessWithData{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Create spending successfully",
			Data:    respMap,
		})

}
