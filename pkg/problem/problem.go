package problem

import "github.com/gofiber/fiber/v2"

type Problem struct {
	Type     string
	Title    string
	Detail   string
	Instance string
	custom   map[string]any
}

func New(c *fiber.Ctx, statusCode int, Type string, title string, detail string, instance string) error {
	return NewProblem(Type, title, detail, instance).Response(c, statusCode)
}

func NewProblem(Type string, title string, detail string, instance string) *Problem {
	return &Problem{
		Type:     Type,
		Title:    title,
		Detail:   detail,
		Instance: instance,
		custom:   make(map[string]any),
	}
}

func (p *Problem) AddCustomFields(fields map[string]any) {
	for k, v := range fields {
		p.custom[k] = v
	}
}

func (p *Problem) Response(c *fiber.Ctx, statusCode ...int) error {
	status := fiber.StatusInternalServerError

	if len(statusCode) > 0 {
		status = statusCode[0]
	}

	m := map[string]any{
		"type":     p.Type,
		"title":    p.Title,
		"detail":   p.Detail,
		"instance": p.Instance,
	}

	for k, v := range p.custom {
		m[k] = v
	}

	return c.Status(status).Type("application/problem+json").JSON(m)
}
