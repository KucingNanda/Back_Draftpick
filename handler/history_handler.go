package handler

import (
	"draftpick-api/model"
	"draftpick-api/repository"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type HistoryRequest struct {
	Role      string `json:"role"`
	PlayerA   string `json:"playerA"`
	PlayerB   string `json:"playerB"`
	DraftData string `json:"draftData"`
}

type DraftPayload struct {
	PlayerA struct {
		FirstHalf struct {
			Picks []PickItem `json:"picks"`
		} `json:"firstHalf"`
		SecondHalf struct {
			Picks []PickItem `json:"picks"`
		} `json:"secondHalf"`
	} `json:"playerA"`
	PlayerB struct {
		FirstHalf struct {
			Picks []PickItem `json:"picks"`
		} `json:"firstHalf"`
		SecondHalf struct {
			Picks []PickItem `json:"picks"`
		} `json:"secondHalf"`
	} `json:"playerB"`
}

type PickItem struct {
	ID          uint   `json:"id"`
	Nama        string `json:"nama"`
	Cost        int    `json:"cost"`
	WeaponLevel int    `json:"weaponLevel"`
	IsC1        bool   `json:"isC1"`
}

func calculateTeamHalfCost(picks []PickItem) int {
	cost := 0
	for _, p := range picks {
		if p.ID == 0 {
			continue
		}
		// Base cost karakter
		cost += p.Cost

		// Penalti senjata B5: R1 (+2), R2 (+3)
		if p.WeaponLevel == 1 {
			cost += 2
		} else if p.WeaponLevel == 2 {
			cost += 3
		}

		// Penalti C1 (+2) pada karakter 5-Star Limited (Cost == 1)
		if p.Cost == 1 && p.IsC1 {
			cost += 2
		}
	}
	return cost
}

// SaveHistory godoc
// @Summary Simpan riwayat pertandingan
// @Description Menyimpan data history match dengan validasi Cost sisi server
// @Tags History
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.History true "Payload history"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/history [post]
func SaveHistory(c *fiber.Ctx) error {
	var req HistoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	if req.Role != "wasit" && req.Role != "peserta" {
		return c.Status(403).JSON(fiber.Map{"error": "Role tidak diizinkan menyimpan riwayat"})
	}

	// === SERVER-SIDE COST VERIFICATION ===
	settings, err := repository.GetSettings()
	if err == nil && settings.MaxCost > 0 && req.DraftData != "" {
		var draftState DraftPayload
		if err := json.Unmarshal([]byte(req.DraftData), &draftState); err == nil {
			costA1 := calculateTeamHalfCost(draftState.PlayerA.FirstHalf.Picks)
			if costA1 > settings.MaxCost {
				return c.Status(400).JSON(fiber.Map{
					"error": fmt.Sprintf("Validasi Gagal: Akumulasi Cost %s (First Half) melebihi batas regulasi (%d / %d poin)", req.PlayerA, costA1, settings.MaxCost),
				})
			}

			costA2 := calculateTeamHalfCost(draftState.PlayerA.SecondHalf.Picks)
			if costA2 > settings.MaxCost {
				return c.Status(400).JSON(fiber.Map{
					"error": fmt.Sprintf("Validasi Gagal: Akumulasi Cost %s (Second Half) melebihi batas regulasi (%d / %d poin)", req.PlayerA, costA2, settings.MaxCost),
				})
			}

			costB1 := calculateTeamHalfCost(draftState.PlayerB.FirstHalf.Picks)
			if costB1 > settings.MaxCost {
				return c.Status(400).JSON(fiber.Map{
					"error": fmt.Sprintf("Validasi Gagal: Akumulasi Cost %s (First Half) melebihi batas regulasi (%d / %d poin)", req.PlayerB, costB1, settings.MaxCost),
				})
			}

			costB2 := calculateTeamHalfCost(draftState.PlayerB.SecondHalf.Picks)
			if costB2 > settings.MaxCost {
				return c.Status(400).JSON(fiber.Map{
					"error": fmt.Sprintf("Validasi Gagal: Akumulasi Cost %s (Second Half) melebihi batas regulasi (%d / %d poin)", req.PlayerB, costB2, settings.MaxCost),
				})
			}
		}
	}
	// ======================================


	limit := int64(10)
	if req.Role == "wasit" {
		limit = 20
	}

	count, _ := repository.GetHistoryCountByRole(req.Role)

	if count >= limit {
		// Limit tercapai, timpa riwayat yang paling lama (UpdatedAt paling tua)
		oldest, err := repository.GetOldestHistoryByRole(req.Role)
		if err == nil {
			oldest.PlayerA = req.PlayerA
			oldest.PlayerB = req.PlayerB
			oldest.DraftData = req.DraftData
			// GORM otomatis memperbarui UpdatedAt
			repository.SaveExistingHistory(&oldest)
			return c.JSON(fiber.Map{"message": "Riwayat lama berhasil ditimpa", "history": oldest})
		}
	}

	// Buat baru jika belum mencapai limit
	newHistory := model.History{
		Role:      req.Role,
		PlayerA:   req.PlayerA,
		PlayerB:   req.PlayerB,
		DraftData: req.DraftData,
	}

	if err := repository.SaveHistory(&newHistory); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan riwayat"})
	}

	return c.JSON(fiber.Map{"message": "Riwayat berhasil disimpan", "history": newHistory})
}

// GetHistories godoc
// @Summary Ambil riwayat pertandingan
// @Description Mengambil daftar seluruh history match
// @Tags History
// @Produce json
// @Success 200 {array} model.History
// @Failure 500 {object} map[string]string
// @Router /api/history [get]
func GetHistories(c *fiber.Ctx) error {
	role := c.Query("role")
	if role == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Parameter role diperlukan"})
	}

	histories, _ := repository.GetHistoriesByRole(role)

	return c.JSON(histories)
}

// DeleteHistory godoc
// @Summary Hapus riwayat pertandingan
// @Description Menghapus history match berdasarkan ID
// @Tags History
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID History"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/history/{id} [delete]
func DeleteHistory(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := repository.DeleteHistory(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus riwayat"})
	}
	return c.JSON(fiber.Map{"message": "Riwayat berhasil dihapus"})
}
