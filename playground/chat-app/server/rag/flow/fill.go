package flow

import (
	"app/server/store"
	"context"
	"log"
)

func DefaultFill(ctx context.Context) *store.Store {
	x := store.New(ctx)
	if x == nil {
		return nil
	}

	// Create bucket with 768 dimensions for gemini-embedding-001.
	// Milvus collection/bucket names allow only letters, numbers, and
	// underscores — no hyphens.
	if err := x.CreateBucket(ctx, "ask_coach", 768, nil); err != nil {
		log.Printf("Failed to create default bucket: %v", err)
		return x
	}

	seed := []struct {
		Content  string
		Metadata map[string]any
	}{
		{
			Content:  "The Italian Game opens 1.e4 e5 2.Nf3 Nc6 3.Bc4, developing quickly and eyeing f7. It's a great starting repertoire for beginners because it teaches rapid development and central control.",
			Metadata: map[string]any{"type": "opening", "topic": "italian-game", "level": "beginner"},
		},
		{
			Content:  "The Sicilian Defense (1.e4 c5) is Black's most popular and aggressive response to 1.e4, fighting for the center asymmetrically instead of mirroring White's pawn.",
			Metadata: map[string]any{"type": "opening", "topic": "sicilian-defense", "level": "intermediate"},
		},
		{
			Content:  "The Queen's Gambit (1.d4 d5 2.c4) offers a pawn to gain central control and open lines. Black is not obligated to keep the pawn; declining with 2...e6 is fully sound.",
			Metadata: map[string]any{"type": "opening", "topic": "queens-gambit", "level": "beginner"},
		},
		{
			Content:  "A fork is a tactic where one piece attacks two or more enemy pieces simultaneously. Knights are especially strong forking pieces because of their unusual movement pattern.",
			Metadata: map[string]any{"type": "tactic", "topic": "fork", "level": "beginner"},
		},
		{
			Content:  "A pin restricts an enemy piece from moving because doing so would expose a more valuable piece behind it, often the king. An absolute pin against the king cannot legally be broken.",
			Metadata: map[string]any{"type": "tactic", "topic": "pin", "level": "beginner"},
		},
		{
			Content:  "A skewer is like a reverse pin: the more valuable piece is attacked first, and when it moves, a less valuable piece behind it is captured.",
			Metadata: map[string]any{"type": "tactic", "topic": "skewer", "level": "intermediate"},
		},
		{
			Content:  "Discovered attacks happen when one piece moves out of the way, revealing an attack from another piece behind it. Discovered checks are especially dangerous since the opponent must respond to the check first.",
			Metadata: map[string]any{"type": "tactic", "topic": "discovered-attack", "level": "intermediate"},
		},
		{
			Content:  "In king and pawn endgames, opposition is critical: the side NOT to move often has the advantage when kings face off. Understanding opposition is essential for converting extra pawns into wins.",
			Metadata: map[string]any{"type": "endgame", "topic": "king-pawn-opposition", "level": "intermediate"},
		},
		{
			Content:  "The Lucena position is a fundamental rook endgame technique for converting a won position: building a 'bridge' with the rook to shield the king from checks while the pawn promotes.",
			Metadata: map[string]any{"type": "endgame", "topic": "lucena-position", "level": "advanced"},
		},
		{
			Content:  "The Philidor position is a key drawing technique in rook endgames where the defending king stays on the promotion rank and the rook checks from behind at the right moment.",
			Metadata: map[string]any{"type": "endgame", "topic": "philidor-position", "level": "advanced"},
		},
		{
			Content:  "Controlling the center early gives pieces more mobility and options. The classical principle is to occupy the center with pawns (e4/d4 or e5/d5) and support it with pieces.",
			Metadata: map[string]any{"type": "strategy", "topic": "center-control", "level": "beginner"},
		},
		{
			Content:  "King safety matters more than material in the opening and middlegame. Castling early, avoiding unnecessary pawn moves near your king, and not delaying development are core safety principles.",
			Metadata: map[string]any{"type": "strategy", "topic": "king-safety", "level": "beginner"},
		},
		{
			Content:  "A common beginner question: 'Why did I lose even though I was up material?' Usually it's from ignoring king safety, hanging pieces to tactics, or entering an endgame with a bad pawn structure despite the material edge.",
			Metadata: map[string]any{"type": "qa", "topic": "material-vs-safety", "level": "beginner"},
		},
		{
			Content:  "A common intermediate question: 'How do I calculate deeper without blundering?' Practice calculating candidate moves in order of forcing-ness (checks, captures, threats first), and always verify your final position before committing.",
			Metadata: map[string]any{"type": "qa", "topic": "calculation-technique", "level": "intermediate"},
		},
	}

	for _, item := range seed {
		embedding, err := embedText(ctx, item.Content)
		if err != nil {
			log.Printf("Failed to embed seed content %q: %v", item.Metadata["topic"], err)
			continue
		}

		if err := x.IngestIfNewWithEmbedding(ctx, "ask_coach", "be_coach", item.Content, item.Metadata, embedding); err != nil {
			log.Printf("Failed to seed %q: %v", item.Metadata["topic"], err)
		}
	}

	return x
}
