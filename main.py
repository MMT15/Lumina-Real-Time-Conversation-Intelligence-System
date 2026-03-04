from fastapi import FastAPI
import ollama

app = FastAPI(title="Lumina Real-Time Intelligence")

@app.get("/")
async def root():
    return {"message": "Lumina Intelligence System is up and running!"}

@app.post("/analyze")
async def analyze_message(text: str):
    # This will be expanded with actual sentiment analysis, ticket creation, and search indexing.
    return {"input": text, "status": "Ready for analysis"}
