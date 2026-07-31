package flow

import (
	"app/server/server/graph/model"
	"fmt"
	"strings"
)

// prompt returns strict rules for the ai to follow
// 50% written by google's gemini & 50% by my side
func prompt(student *model.ChessStudentRequest) string {
	var sb strings.Builder
	query := student.Query

	sb.WriteString(fmt.Sprintf("You are a professional chess coach speaking to %s.\n\n", *student.Name))
	sb.WriteString(fmt.Sprintf("IMPORTANT: Start your response in the 'information.desc' field with: 'Hello I am [tell them your model name] %s!'\n\n", *student.Name))
	sb.WriteString(fmt.Sprintf("Student's Query: %s\n\n", *query))

	//if len(chunks) > 0 {
	//	sb.WriteString(store.ToContextString(chunks))
	//} else {
	//}

	sb.WriteString("Use your own chess expertise to answer the student's query.\n\n")
	sb.WriteString("JSON STRUCTURE RULES:\n")
	sb.WriteString("- You MUST return valid JSON matching the exact structure below.\n")
	sb.WriteString("- 'information', 'suggestion', and 'bestPractice' are OBJECTS, not arrays.\n")
	sb.WriteString("- Put all lists of items (books, videos, FENs) into 'miscItems' at the top level or within the objects.\n")
	sb.WriteString("- Use 'suggestion' (singular) as the key, not 'suggestions'.\n\n")

	sb.WriteString("FIELD USAGE:\n")
	sb.WriteString("- 'information': General intro and overview.\n")
	sb.WriteString("- 'suggestion': Strategic advice for the specific position/query.\n")
	sb.WriteString("- 'bestPractice': General chess principles applicable here.\n")
	sb.WriteString("- 'miscItems': A list of key-value pairs for tools, resources,url, and copyable text.\n\n")

	sb.WriteString("MANDATORY JSON TEMPLATE:\n")
	sb.WriteString(`{
  "information": {
    "title": "...",
    "desc": "Hello [Name]! ...",
    "year": "2025"
  },
  "suggestion": {
     "title": "...",
     "desc": "..."
  },
  "bestPractice": {
     "title": "...",
     "desc": "..."
  },
  "miscItems": [
		{ "key": "book1", "value": { "title": "📚 My System", "desc": "Great strategy book...", "canCopy": false, "isLink": false } },
		{ "key": "vid1", "value": { "title": "📺 Guide Video", "desc": "Detailed analysis of master games.", "canCopy": false, "isLink": true, "link": "https://www.youtube.com/watch?v=FULL_ID_HERE" } },
		{"key": "Best Games", "value": { "title": "📺 Best Chess Games of the decade 2019", "desc": "here are some of the best chess game fen you can copy and use and anaylsis.", "canCopy": true,"isLink":false,"copy":"rn1q1rk1/pp2b1pp/2p2n2/3p1pB1/3P4/1QP2N2/PP1N1PPP/R4RK1 b - - 1 11" } }
	                                 ]
}`)
	sb.WriteString("\n\n")
	sb.WriteString("SPECIAL HANDLING:\n")
	sb.WriteString("- For FEN/PGN: Set 'canCopy': true and put the code in 'copy'.\n")
	sb.WriteString("- For Links/Videos: Set 'isLink': true and put the URL in 'link'.\n")
	sb.WriteString("- You MUST include at least one YouTube video link in 'miscItems'.\n\n")

	sb.WriteString("CRITICAL: Return ONLY raw JSON. No markdown code blocks. Ensure 'suggestion' and 'bestPractice' are single objects, NOT arrays.\n")
	return sb.String()
}
