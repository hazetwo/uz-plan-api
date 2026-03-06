from pathlib import Path
from typing import ClassVar

from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    API_V1_STR: str = "/api/v1"
    PROJECT_NAME: str = "uz-plan-api"
    SCHEDULE_LINK: str = "https://plan.uz.zgora.pl/grupy_plan.php?ID="
    SUPPORTED_MAJOR: str = (
        "https://plan.uz.zgora.pl/grupy_lista_grup_kierunku.php?ID=401"
    )
    GROUPS_FILE: ClassVar[Path] = (
        Path(__file__).resolve().parents[2] / "groups.json"
    )
    LONGEST_ID: int = 5


settings = Settings()
