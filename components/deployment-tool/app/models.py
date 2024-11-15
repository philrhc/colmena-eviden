from pydantic import BaseModel

class InputData(BaseModel):
    base_directory: str
    repo_url: str