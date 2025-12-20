package cap

import (
	"fmt"
)

// Translator handles translation of CAP message content
type Translator struct {
	// In production, this would integrate with a translation service
	// For PoC, we'll use simple placeholder logic
}

// NewTranslator creates a new translator
func NewTranslator() *Translator {
	return &Translator{}
}

// Translate translates text to target language
// This is a placeholder implementation - in production would use a real translation service
func (t *Translator) Translate(text string, sourceLang string, targetLang string) (string, error) {
	// Placeholder: return original text with language prefix
	// In production, this would call a translation API
	return fmt.Sprintf("[%s] %s", targetLang, text), nil
}

// TranslateCAPInfo translates CAP Info block to target language
func (t *Translator) TranslateCAPInfo(info Info, targetLang string) (*Info, error) {
	translated := Info{
		Language:    targetLang,
		Category:    info.Category, // Categories don't need translation
		Event:       info.Event,     // Event types don't need translation (standard codes)
		Urgency:     info.Urgency,   // Urgency levels don't need translation (standard codes)
		Severity:    info.Severity,  // Severity levels don't need translation (standard codes)
		Certainty:   info.Certainty, // Certainty levels don't need translation (standard codes)
		Contact:     info.Contact,   // Contact info doesn't need translation
		Expires:     info.Expires,   // Timestamps don't need translation
	}
	
	// Translate language-specific fields
	var err error
	translated.Headline, err = t.Translate(info.Headline, info.Language, targetLang)
	if err != nil {
		return nil, fmt.Errorf("failed to translate headline: %w", err)
	}
	
	translated.Description, err = t.Translate(info.Description, info.Language, targetLang)
	if err != nil {
		return nil, fmt.Errorf("failed to translate description: %w", err)
	}
	
	translated.Instruction, err = t.Translate(info.Instruction, info.Language, targetLang)
	if err != nil {
		return nil, fmt.Errorf("failed to translate instruction: %w", err)
	}
	
	return &translated, nil
}

// SupportedLanguages returns list of supported languages
func (t *Translator) SupportedLanguages() []string {
	return []string{"zh-TW", "en", "zh-CN", "ja", "ko"}
}

