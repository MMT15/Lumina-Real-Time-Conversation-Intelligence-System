import ollama

def analyze_sentiment(text):
    response = ollama.chat(model='llama3.2:3b', messages=[
      {
        'role': 'user',
        'content': f"Analyze the sentiment of this message: '{text}'. Return ONLY a JSON with keys: 'sentiment' (string) and 'score' (0-1).",
      },
    ])
    print(response['message']['content'])

analyze_sentiment("I am very frustrated, my order hasn't arrived yet!")
