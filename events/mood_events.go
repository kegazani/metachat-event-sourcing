package events

// MoodAnalyzedPayload represents the payload for MoodAnalyzed event
type MoodAnalyzedPayload struct {
	DiaryEntryID string                 `json:"diary_entry_id"`
	Emotions     map[string]float64    `json:"emotions"`
	Dominant     string                 `json:"dominant_emotion"`
	Valence      float64                `json:"valence"`
	Arousal      float64                `json:"arousal"`
	Confidence   float64                `json:"confidence"`
	Keywords     []string               `json:"keywords,omitempty"`
	Topics       []string               `json:"topics,omitempty"`
}

// DailyMoodAggregatedPayload represents the payload for DailyMoodAggregated event
type DailyMoodAggregatedPayload struct {
	UserID       string                 `json:"user_id"`
	Date         string                 `json:"date"`
	Emotions     map[string]float64    `json:"emotions"`
	Dominant     string                 `json:"dominant_emotion"`
	Valence      float64                `json:"valence"`
	Arousal      float64                `json:"arousal"`
	EntryCount   int                    `json:"entry_count"`
	TokenCount   int                    `json:"token_count"`
	Volatility   float64                `json:"volatility"`
	Keywords     []string               `json:"keywords,omitempty"`
	Topics       []string               `json:"topics,omitempty"`
}

// WeeklyMoodAggregatedPayload represents the payload for WeeklyMoodAggregated event
type WeeklyMoodAggregatedPayload struct {
	UserID       string                 `json:"user_id"`
	Week         int                    `json:"week"`
	Year         int                    `json:"year"`
	Emotions     map[string]float64    `json:"emotions"`
	Dominant     string                 `json:"dominant_emotion"`
	Valence      float64                `json:"valence"`
	Arousal      float64                `json:"arousal"`
	EntryCount   int                    `json:"entry_count"`
	TokenCount   int                    `json:"token_count"`
	Volatility   float64                `json:"volatility"`
	Trend        string                 `json:"trend"` // "improving", "declining", "stable"
	Keywords     []string               `json:"keywords,omitempty"`
	Topics       []string               `json:"topics,omitempty"`
}

// MonthlyMoodAggregatedPayload represents the payload for MonthlyMoodAggregated event
type MonthlyMoodAggregatedPayload struct {
	UserID       string                 `json:"user_id"`
	Month        int                    `json:"month"`
	Year         int                    `json:"year"`
	Emotions     map[string]float64    `json:"emotions"`
	Dominant     string                 `json:"dominant_emotion"`
	Valence      float64                `json:"valence"`
	Arousal      float64                `json:"arousal"`
	EntryCount   int                    `json:"entry_count"`
	TokenCount   int                    `json:"token_count"`
	Volatility   float64                `json:"volatility"`
	Trend        string                 `json:"trend"` // "improving", "declining", "stable"
	Keywords     []string               `json:"keywords,omitempty"`
	Topics       []string               `json:"topics,omitempty"`
}

// UserPortraitUpdatedPayload represents the payload for UserPortraitUpdated event
type UserPortraitUpdatedPayload struct {
	UserID           string                 `json:"user_id"`
	EmotionalProfile EmotionalProfile      `json:"emotional_profile"`
	BehavioralProfile BehavioralProfile    `json:"behavioral_profile"`
	ThematicProfile  ThematicProfile       `json:"thematic_profile"`
	ArchetypeProfile ArchetypeProfile      `json:"archetype_profile"`
	Modalities       []UserModality        `json:"modalities"`
	LastUpdated      string                `json:"last_updated"`
}

// EmotionalProfile represents the emotional characteristics of a user
type EmotionalProfile struct {
	BaseEmotions     map[string]float64 `json:"base_emotions"`
	Valence          float64            `json:"valence"`
	Arousal          float64            `json:"arousal"`
	EmotionalRange   float64            `json:"emotional_range"`
	Stability        float64            `json:"stability"`
	Reactivity       float64            `json:"reactivity"`
}

// BehavioralProfile represents the behavioral patterns of a user
type BehavioralProfile struct {
	EntryFrequency   map[string]int    `json:"entry_frequency"` // entries by day of week, time of day
	AverageEntryLength int             `json:"average_entry_length"`
	SessionPatterns  map[string]int    `json:"session_patterns"`
	ActivityTimeline map[string]int    `json:"activity_timeline"`
}

// ThematicProfile represents the thematic preferences of a user
type ThematicProfile struct {
	TopTopics        []TopicWeight     `json:"top_topics"`
	Keywords         []string          `json:"keywords"`
	Interests        []string          `json:"interests"`
	WritingStyle     WritingStyle      `json:"writing_style"`
}

// TopicWeight represents a topic with its weight
type TopicWeight struct {
	Topic  string  `json:"topic"`
	Weight float64 `json:"weight"`
}

// WritingStyle represents the writing style characteristics
type WritingStyle struct {
	AverageSentenceLength float64   `json:"average_sentence_length"`
	VocabularyComplexity  float64   `json:"vocabulary_complexity"`
	Emotiveness          float64   `json:"emotiveness"`
	Formality            float64   `json:"formality"`
	CommonWords          []string  `json:"common_words"`
}

// ArchetypeProfile represents the archetype characteristics of a user
type ArchetypeProfile struct {
	PrimaryArchetype   Archetype        `json:"primary_archetype"`
	SecondaryArchetypes []Archetype     `json:"secondary_archetypes"`
	ArchetypeScores    map[string]float64 `json:"archetype_scores"`
}

// Archetype represents a psychological archetype
type Archetype struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Score       float64 `json:"score"`
	Traits      []string `json:"traits"`
}