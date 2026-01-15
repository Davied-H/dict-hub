package vocabulary

import (
	"time"

	"dict-hub/internal/model"
)

// 艾宾浩斯复习间隔（天数）
// Level 0: 新词，立即复习
// Level 1: 1天后
// Level 2: 2天后
// Level 3: 4天后
// Level 4: 7天后
// Level 5: 15天后
// Level 6: 30天后
// Level 7: 60天后（已掌握）
var ReviewIntervals = []time.Duration{
	0,                    // Level 0: 立即
	24 * time.Hour,       // Level 1: 1天
	2 * 24 * time.Hour,   // Level 2: 2天
	4 * 24 * time.Hour,   // Level 3: 4天
	7 * 24 * time.Hour,   // Level 4: 7天
	15 * 24 * time.Hour,  // Level 5: 15天
	30 * 24 * time.Hour,  // Level 6: 30天
	60 * 24 * time.Hour,  // Level 7: 60天
}

const MaxLevel = 7

// CalculateNextReview 根据复习结果计算下次复习时间和新等级
func CalculateNextReview(currentLevel int, result model.ReviewResult) (newLevel int, nextReview time.Time) {
	now := time.Now()

	switch result {
	case model.ReviewResultForgot:
		// 完全忘记：降到0级，立即需要复习
		newLevel = 0
		nextReview = now

	case model.ReviewResultHard:
		// 困难：降一级（最低为0）
		newLevel = currentLevel - 1
		if newLevel < 0 {
			newLevel = 0
		}
		nextReview = now.Add(ReviewIntervals[newLevel])

	case model.ReviewResultGood:
		// 记得：升一级
		newLevel = currentLevel + 1
		if newLevel > MaxLevel {
			newLevel = MaxLevel
		}
		nextReview = now.Add(ReviewIntervals[newLevel])

	case model.ReviewResultEasy:
		// 非常容易：升两级
		newLevel = currentLevel + 2
		if newLevel > MaxLevel {
			newLevel = MaxLevel
		}
		nextReview = now.Add(ReviewIntervals[newLevel])

	default:
		// 默认保持当前等级
		newLevel = currentLevel
		nextReview = now.Add(ReviewIntervals[currentLevel])
	}

	return newLevel, nextReview
}

// GetLevelDescription 获取等级描述
func GetLevelDescription(level int) string {
	descriptions := []string{
		"新词",     // 0
		"初学",     // 1
		"学习中",   // 2
		"熟悉中",   // 3
		"基本掌握", // 4
		"良好掌握", // 5
		"熟练掌握", // 6
		"已掌握",   // 7
	}
	if level < 0 || level >= len(descriptions) {
		return "未知"
	}
	return descriptions[level]
}

// GetIntervalDescription 获取间隔描述
func GetIntervalDescription(level int) string {
	if level < 0 || level > MaxLevel {
		return "未知"
	}
	intervals := []string{
		"立即",   // 0
		"1天",    // 1
		"2天",    // 2
		"4天",    // 3
		"7天",    // 4
		"15天",   // 5
		"30天",   // 6
		"60天",   // 7
	}
	return intervals[level]
}
