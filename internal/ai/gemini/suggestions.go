package gemini

import (
	"encoding/json"

	"github.com/google/generative-ai-go/genai"
	"github.com/oyen-bright/goFundIt/internal/models"
)

// TODO; Add location ?
func (c *geminiClient) GenerateActivitySuggestions(prompt string) ([]models.ActivitySuggestion, error) {
	c.model.SetTemperature(0.4)
	c.model.SetTopK(40)
	c.model.SetTopP(0.95)
	c.model.SetMaxOutputTokens(8192)
	c.model.ResponseMIMEType = "application/json"
	c.model.ResponseSchema = &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"activityName": {
					Type: genai.TypeString,
				},
				"estimatedCost": {
					Type: genai.TypeString,
				},
			},
		},
	}

	session := c.model.StartChat()
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("I want an AI that will suggest activities for a provided title and description or holiday trip with real-life data. I want the activity names and estimated costs."),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("Okay, let's break down how we can approach building an AI that suggests activities for a provided title or holiday trip, including estimated costs based on real-world data. This is a complex project, but we can outline the core components and how they would work together.\n\n**Core Components of the AI System:**\n\n1.  **Natural Language Processing (NLP) Module:**\n    *   **Purpose:** To understand the user's input (title or holiday trip description).\n    *   **Tasks:**\n        *   **Entity Recognition:** Identify key entities like locations, dates, interests (e.g., \"hiking,\" \"museums,\" \"beach\"), and the number of people.\n        *   **Intent Classification:** Determine the user's goal (e.g., \"find activities,\" \"plan a trip,\" \"find things to do near X\").\n        *   **Sentiment Analysis:** Understand the user's mood/preferences (if provided).\n\n2.  **Knowledge Base:**\n    *   **Purpose:** A structured database of activities, locations, and costs.\n    *   **Data Sources:** This is the most challenging part. We'd need to gather and regularly update data from various sources, such as:\n        *   **Activity APIs:**\n            *   **GetYourGuide API:** Offers a wide range of tours and activities globally.\n            *   **Viator API:** Another significant source for tours and attractions.\n            *   **Eventbrite API:** For local events and happenings.\n            *   **Google Places API:** For finding local attractions, points of interest, and reviews.\n        *   **Travel Aggregators:**\n            *   **Expedia, Booking.com:**  May provide some activity-related info.\n        *   **Tourism Boards:** Local tourism websites often list activities and their prices.\n        *   **Crowd-Sourced Data:**  Web scraping (with care and ethical considerations) from sites like TripAdvisor, Yelp, and local review platforms.\n        *   **Open Data Sets:** If any exist for tourism-related activity pricing and info.\n    *   **Data Structure:**\n        *   **Activity Information:** Name, description, location, keywords, duration, required items, images, average rating, pricing (including ranges), seasonality.\n        *   **Location Information:** Geo-coordinates, type (city, landmark, natural area), climate, accessibility.\n        *   **User Reviews/Ratings:** To build a \"popular\" filter.\n        *   **Cost Data:** May require aggregation and averaging across sources. This might be a 'per person' estimation, or a combined family/group option.\n\n3.  **Matching and Ranking Algorithm:**\n    *   **Purpose:** To find activities that match the user's input and rank them by relevance and suitability.\n    *   **Steps:**\n        1.  **Filtering:** Based on location, date, and user-specified criteria (e.g., \"outdoors,\" \"family-friendly\").\n        2.  **Keyword Matching:** Compare extracted keywords from user input with keywords associated with activities.\n        3.  **Relevance Ranking:**  Combine various signals (keyword match, user rating, price, location proximity) to rank results.\n        4.  **Diversity:**  Ensure the results aren't all too similar - offer a variety if possible.\n\n4.  **Cost Estimation Module:**\n    *   **Purpose:** To provide estimated costs for activities and possibly associated expenses.\n    *   **Tasks:**\n        *   **Price Lookup:** Retrieve the price (or average price) from the database for the selected activity.\n        *   **Calculate Total:** Multiply the per-person price by the number of people.\n        *   **Additional Costs:** Consider other expenses (if possible):\n            *   **Transportation:** Estimated transportation costs to and from the activity location (if provided within the activity info).\n            *   **Equipment Rental:** If specified.\n            *   **Food/Drink:** Rough approximations, based on the activity and area.\n        *   **Currency Conversion:** If needed.\n    *   **Important Note:** Cost estimation will always have an element of uncertainty. We can provide \"low,\" \"average,\" and \"high\" estimates.\n\n5.  **Output Generator:**\n    *   **Purpose:** To present the results to the user in an understandable way.\n    *   **Tasks:**\n        *   **Activity List:** Provide a formatted list of activities, including descriptions, images, and estimated costs.\n        *   **Filtering and Sorting:** Allow users to filter results by price, rating, or other criteria.\n        *   **Interactive Map:** If needed, present a map with activity locations.\n\n**Implementation:**\n\n*   **Programming Languages:** Python is a strong choice due to its libraries for NLP, data analysis, and web scraping.\n*   **Libraries:**\n    *   **NLP:** NLTK, SpaCy, Transformers (Hugging Face)\n    *   **Data Handling:** Pandas, NumPy\n    *   **Web Scraping:** Beautiful Soup, Scrapy\n    *   **APIs:** Requests (to access external APIs)\n    *   **Database:** SQLite, PostgreSQL, or a vector database\n*   **Cloud Services (optional):** AWS, Google Cloud, Azure could be used for hosting and scaling the AI system.\n*   **User Interface:** It could be a simple web application or a mobile app.\n\n**Challenges:**\n\n*   **Data Quality:** Gathering and maintaining clean and accurate cost data is challenging.\n*   **Scalability:** As the database grows, the system needs to handle the increased load.\n*   **Dynamic Prices:** The prices of activities can fluctuate and be seasonal, so it will need to be updated regularly.\n*   **Geographic Coverage:** Covering the whole world with reliable data is very ambitious. It might be best to focus on specific regions.\n*   **User Feedback:** Incorporating user feedback to improve recommendations.\n*   **Cost Accuracy:** Prices can vary a lot so we may need to provide estimates, ranges and disclaimers.\n\n**Simple Example (Illustrative):**\n\nLet's say the user enters \"Weekend trip to New York City with kids, museum and park activities\".\n\n1.  **NLP:** Extracts \"New York City,\" \"kids,\" \"museum,\" \"park,\" and \"weekend\"\n2.  **Knowledge Base:** The database contains details of various New York City museums, parks and kid-friendly activities with prices.\n3.  **Matching/Ranking:** Results: The system filters for kid-friendly locations in NYC,  matches \"museum\" and \"park\" with activities such as The Met Museum, Central Park bike rental, The Bronx Zoo. It ranks results based on popularity and keywords, as well as 'suitable for kids' tag.\n4.  **Cost Estimation:** The system estimates cost for museum entry tickets, the bike rental and zoo tickets, and perhaps a budget for food and snacks.\n5.  **Output:** Lists these activity options with photos, descriptions, and estimated cost ranges.\n\n**Next Steps for You:**\n\n1.  **Start Small:** Don't try to cover everything initially. Focus on a single region or activity type.\n2.  **Data Gathering:** Explore the various data sources, pick some to work with, and start building a database.\n3.  **Choose a Stack:** Select your tech, and start with a basic pipeline.\n4.  **Iterate:** Build, test, and then improve based on results.\n\nThis is a very complex project, but if you take it step-by-step you can create a useful and effective tool.  Do you want to start with a specific step or dive more into any of the components listed?\n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("I will provide the current currency (e.g., NGN, USD, EUR) as input.\nThe estimated cost should be displayed in the same currency."),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n"),
				genai.Text("[\n  {\n    \"activityName\": \"Example Activity 1\",\n    \"estimatedCost\": \"150 USD\"\n  },\n    {\n    \"activityName\": \"Example Activity 2\",\n    \"estimatedCost\": \"12000 NGN\"\n  },\n    {\n    \"activityName\": \"Example Activity 3\",\n    \"estimatedCost\": \"75 EUR\"\n  }\n]"),
				genai.Text("\n```"),
			},
		},
	}

	resp, err := session.SendMessage(*c.context, genai.Text(prompt))

	if err != nil {
		return nil, err
	}

	var suggestions []models.ActivitySuggestion
	for _, part := range resp.Candidates[0].Content.Parts {
		if s, ok := part.(genai.Text); ok {

			var parsedSuggestions []models.ActivitySuggestion
			if err := json.Unmarshal([]byte(s), &parsedSuggestions); err != nil {
				return nil, err
			}
			suggestions = append(suggestions, parsedSuggestions...)

		}
	}
	return suggestions, nil
}
