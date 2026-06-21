package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"wxcloudrun-golang/internal/model"
	"wxcloudrun-golang/internal/repository"
	"wxcloudrun-golang/internal/types"
	"wxcloudrun-golang/internal/util"
)

type ParseRecordInput struct {
	FamilyID          uint   `json:"family_id"`
	Text              string `json:"text"`
	CurrentBabyID     uint   `json:"current_baby_id"`
	TimezoneOffsetMin int    `json:"timezone_offset_min"`
}

type ChatInput struct {
	FamilyID          uint   `json:"family_id"`
	BabyID            uint   `json:"baby_id"`
	Text              string `json:"text"`
	TimezoneOffsetMin int    `json:"timezone_offset_min"`
}

func ParseRecord(openID string, input ParseRecordInput) (map[string]interface{}, error) {
	if _, _, err := RequireFamilyMember(openID, input.FamilyID); err != nil {
		return nil, err
	}

	babies, err := repository.ListBabiesByFamilyID(input.FamilyID)
	if err != nil {
		return nil, err
	}

	if len(babies) == 0 {
		return map[string]interface{}{
			"type":         "record_parse",
			"need_confirm": true,
			"records":      []types.RecordPayload{},
			"warnings":     []string{"还没有宝宝档案，请先添加宝宝。"},
		}, nil
	}

	text := strings.TrimSpace(input.Text)
	baby := pickBaby(babies, text, input.CurrentBabyID)
	warnings := []string{}

	clock := parseClock(text)
	actionTime := time.Now().UTC().Format(time.RFC3339)
	if clock != "" {
		actionTime = util.LocalActionTimeFromClock(clock, input.TimezoneOffsetMin).UTC().Format(time.RFC3339)
	} else {
		warnings = append(warnings, "没有识别到明确时间，先按当前时间处理。")
	}

	records := []types.RecordPayload{}

	if amountMatch := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*ml`).FindStringSubmatch(text); len(amountMatch) > 1 && strings.ContainsAny(text, "奶喝") {
		records = append(records, types.RecordPayload{
			BabyID:          baby.ID,
			BabyName:        displayBabyName(&baby),
			ActionType:      "feed",
			ActionTime:      actionTime,
			Summary:         fmt.Sprintf("%s 配方奶 %sml", displayBabyName(&baby), amountMatch[1]),
			Data:            map[string]interface{}{"amount_ml": toInt(amountMatch[1]), "feed_type": feedTypeFromText(text)},
			Source:          "ai",
			ClientRequestID: newClientRequestID(),
		})
	}

	hasPee := strings.Contains(text, "尿")
	hasPoop := strings.Contains(text, "便") || strings.Contains(text, "拉")
	if hasPee || hasPoop {
		records = append(records, types.RecordPayload{
			BabyID:          baby.ID,
			BabyName:        displayBabyName(&baby),
			ActionType:      "diaper",
			ActionTime:      actionTime,
			Summary:         fmt.Sprintf("%s %s%s", displayBabyName(&baby), ifThen(hasPee, "尿", ""), ifThen(hasPoop, "便", "")),
			Data:            map[string]interface{}{"pee": hasPee, "poop": hasPoop, "poop_color": poopColorFromText(text)},
			Source:          "ai",
			ClientRequestID: newClientRequestID(),
		})
	}

	if strings.Contains(text, "睡") {
		if duration := parseDurationMin(text); duration > 0 {
			records = append(records, types.RecordPayload{
				BabyID:          baby.ID,
				BabyName:        displayBabyName(&baby),
				ActionType:      "sleep",
				ActionTime:      actionTime,
				Summary:         fmt.Sprintf("%s 睡了 %d 分钟", displayBabyName(&baby), duration),
				Data:            map[string]interface{}{"duration_min": duration},
				Source:          "ai",
				ClientRequestID: newClientRequestID(),
			})
		}
	}

	if strings.Contains(text, "体重") {
		if weight := parseWeightG(text); weight > 0 {
			records = append(records, types.RecordPayload{
				BabyID:          baby.ID,
				BabyName:        displayBabyName(&baby),
				ActionType:      "weight",
				ActionTime:      actionTime,
				Summary:         fmt.Sprintf("%s 体重 %dg", displayBabyName(&baby), weight),
				Data:            map[string]interface{}{"weight_g": weight},
				Source:          "ai",
				ClientRequestID: newClientRequestID(),
			})
		}
	}

	if strings.Contains(text, "第一次") || strings.Contains(text, "翻身") || strings.Contains(text, "抬头") || strings.Contains(text, "笑出声") {
		records = append(records, types.RecordPayload{
			BabyID:          baby.ID,
			BabyName:        displayBabyName(&baby),
			ActionType:      "milestone",
			ActionTime:      actionTime,
			Summary:         fmt.Sprintf("%s %s", displayBabyName(&baby), text),
			Data:            map[string]interface{}{"title": text, "content": ""},
			Source:          "ai",
			ClientRequestID: newClientRequestID(),
		})
	}

	if strings.Contains(text, "AD") || strings.Contains(text, "用药") || strings.Contains(text, "药") {
		records = append(records, types.RecordPayload{
			BabyID:          baby.ID,
			BabyName:        displayBabyName(&baby),
			ActionType:      "medicine",
			ActionTime:      actionTime,
			Summary:         fmt.Sprintf("%s 用药 %s", displayBabyName(&baby), text),
			Data:            map[string]interface{}{"medicine_name": text},
			Source:          "ai",
			ClientRequestID: newClientRequestID(),
		})
	}

	return map[string]interface{}{
		"type":         "record_parse",
		"need_confirm": true,
		"records":      records,
		"warnings":     warnings,
	}, nil
}

func Chat(openID string, input ChatInput) (map[string]interface{}, error) {
	if _, _, err := RequireFamilyMember(openID, input.FamilyID); err != nil {
		return nil, err
	}

	text := strings.TrimSpace(input.Text)
	if strings.Contains(text, "发烧") || strings.Contains(text, "黄疸") || strings.Contains(text, "异常") || strings.Contains(text, "用药建议") {
		return map[string]interface{}{
			"answer":          "这个问题涉及宝宝个体情况，我不能替代医生判断。我可以帮你整理最近的记录，方便你咨询医生。如果宝宝出现精神反应差、吃奶明显减少、呼吸异常、发热、持续呕吐等情况，建议及时联系医生或就医。",
			"related_actions": []ActionView{},
		}, nil
	}

	start, end, _ := util.TodayRangeByOffset(input.TimezoneOffsetMin)
	actions, err := repository.ListActions(repository.ListActionsInput{
		FamilyID:  input.FamilyID,
		BabyID:    input.BabyID,
		StartTime: &start,
		EndTime:   &end,
		Limit:     200,
		Offset:    0,
	})
	if err != nil {
		return nil, err
	}
	if len(actions) == 0 {
		return map[string]interface{}{
			"answer":          "目前没有相关记录。",
			"related_actions": []ActionView{},
		}, nil
	}

	related := make([]ActionView, 0, min(5, len(actions)))
	for i := 0; i < min(5, len(actions)); i++ {
		related = append(related, buildActionView(actions[i]))
	}

	switch {
	case strings.Contains(text, "喝了多少") || strings.Contains(text, "奶量"):
		total := 0
		count := 0
		for _, action := range actions {
			if action.ActionType == "feed" {
				count++
				total += toInt(util.MustUnmarshal(action.DataJSON)["amount_ml"])
			}
		}
		return map[string]interface{}{
			"answer":          fmt.Sprintf("今天目前记录了 %d 次喂养，总量 %dml。", count, total),
			"related_actions": related,
		}, nil
	case strings.Contains(text, "尿了几次"):
		count := 0
		for _, action := range actions {
			if action.ActionType == "diaper" && toBool(util.MustUnmarshal(action.DataJSON)["pee"]) {
				count++
			}
		}
		return map[string]interface{}{
			"answer":          fmt.Sprintf("今天目前记录了 %d 次尿尿。", count),
			"related_actions": related,
		}, nil
	case strings.Contains(text, "睡了多久") || strings.Contains(text, "睡眠"):
		total := 0
		for _, action := range actions {
			if action.ActionType == "sleep" {
				total += toInt(util.MustUnmarshal(action.DataJSON)["duration_min"])
			}
		}
		return map[string]interface{}{
			"answer":          fmt.Sprintf("今天目前累计睡了 %d 分钟。", total),
			"related_actions": related,
		}, nil
	default:
		summaries := []string{}
		for _, action := range actions[:min(5, len(actions))] {
			summaries = append(summaries, action.Summary)
		}
		return map[string]interface{}{
			"answer":          "今天的重点记录有：" + strings.Join(summaries, "；") + "。",
			"related_actions": related,
		}, nil
	}
}

func pickBaby(babies []model.Baby, text string, currentBabyID uint) model.Baby {
	for _, baby := range babies {
		if baby.Nickname != "" && strings.Contains(text, baby.Nickname) {
			return baby
		}
		if strings.Contains(text, baby.Name) {
			return baby
		}
	}
	for _, baby := range babies {
		if baby.ID == currentBabyID {
			return baby
		}
	}
	return babies[0]
}

func parseClock(text string) string {
	match := regexp.MustCompile(`(\d{1,2})\s*点(?:半|(\d{1,2})分?)?`).FindStringSubmatch(text)
	if len(match) == 0 {
		return ""
	}
	minute := "00"
	if strings.Contains(match[0], "半") {
		minute = "30"
	} else if len(match) > 2 && match[2] != "" {
		minute = fmt.Sprintf("%02s", match[2])
	}
	return fmt.Sprintf("%02s:%s", match[1], minute)
}

func parseDurationMin(text string) int {
	match := regexp.MustCompile(`(\d+)\s*(分钟|小时)`).FindStringSubmatch(text)
	if len(match) < 3 {
		return 0
	}
	value := toInt(match[1])
	if match[2] == "小时" {
		return value * 60
	}
	return value
}

func parseWeightG(text string) int {
	match := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(kg|公斤|g|克)`).FindStringSubmatch(text)
	if len(match) < 3 {
		return 0
	}
	if match[2] == "kg" || match[2] == "公斤" {
		return int(toFloat(match[1]) * 1000)
	}
	return toInt(match[1])
}

func feedTypeFromText(text string) string {
	if strings.Contains(text, "母乳") {
		return "breast"
	}
	return "formula"
}

func poopColorFromText(text string) string {
	for _, color := range []string{"黄色", "绿色", "棕色", "金黄", "褐色"} {
		if strings.Contains(text, color) {
			return color
		}
	}
	return ""
}

func newClientRequestID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix()%10000)
}

func ifThen(cond bool, yes, no string) string {
	if cond {
		return yes
	}
	return no
}

func toFloat(value string) float64 {
	var result float64
	fmt.Sscanf(value, "%f", &result)
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
