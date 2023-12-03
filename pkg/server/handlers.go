package server

import (
	"github.com/1x-eng/mbat/pkg/job"
	"github.com/1x-eng/mbat/pkg/microbatcher"
	"github.com/gofiber/fiber/v2"
)

func jobSubmitHandler(batcher *microbatcher.MicroBatcher) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jData := new(struct {
			Data string `json:"data"`
		})

		if err := c.BodyParser(jData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		j := job.NewJob(jData.Data)
		resultChan := batcher.Submit(j)

		result := <-resultChan
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
		}

		return c.JSON(result)
	}
}
