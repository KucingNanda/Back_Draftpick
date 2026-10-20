package handler

import (
	"draftpick-api/model"
	"draftpick-api/repository"

	"github.com/gofiber/fiber/v2"
)

type UserWithPesertaResponse struct {
	ID       uint           `json:"id"`
	Username string         `json:"username"`
	Role     string         `json:"role"`
	Peserta  *model.Peserta `json:"peserta,omitempty"`
}

// GetUsers godoc
// @Summary Ambil daftar seluruh akun pengguna
// @Description Mengambil semua data pengguna beserta data profil peserta miliknya (khusus wasit)
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} UserWithPesertaResponse
// @Failure 500 {object} map[string]string
// @Router /api/users [get]
func GetUsers(c *fiber.Ctx) error {
	roleQuery := c.Query("role")
	var users []model.User
	var err error

	if roleQuery != "" {
		users, err = repository.GetUsersByRole(roleQuery)
	} else {
		users, err = repository.GetAllUsers()
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data pengguna"})
	}

	var response []UserWithPesertaResponse
	for _, u := range users {
		peserta, _ := repository.GetPesertaByUserID(u.ID)
		response = append(response, UserWithPesertaResponse{
			ID:       u.ID,
			Username: u.Username,
			Role:     u.Role,
			Peserta:  peserta,
		})
	}

	return c.JSON(response)
}

// UpdateUserRole godoc
// @Summary Ubah role pengguna (Promote / Demote)
// @Description Mengubah role user (misal dari member menjadi peserta atau sebaliknya)
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID User"
// @Param request body model.UpdateUserRoleRequest true "Payload role baru"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id}/role [put]
func UpdateUserRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateUserRoleRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format payload tidak valid"})
	}

	if req.Role != "member" && req.Role != "peserta" && req.Role != "wasit" {
		return c.Status(400).JSON(fiber.Map{"error": "Role harus salah satu dari: member, peserta, wasit"})
	}

	if err := repository.UpdateUserRole(id, req.Role); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengubah role pengguna"})
	}

	return c.JSON(fiber.Map{
		"message": "Role pengguna berhasil diperbarui!",
		"role":    req.Role,
	})
}

