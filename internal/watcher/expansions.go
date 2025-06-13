package watcher

import "github.com/gofiber/fiber/v2"

func (w *Watcher) GetExpansions(c *fiber.Ctx) error {
	expansions, err := w.cardtraderAdapter.GetExpansions(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(expansions)
}
