import json

import httpx

from app.config.settings import settings
from app.models import Group


def get_url_by_id(id: str) -> httpx.URL | None:
    if not id.isdigit() or len(id) > settings.LONGEST_ID:
        return None

    with settings.GROUPS_FILE.open("r", encoding="utf-8") as file:
        data = json.load(file)

    groups = [Group.model_validate(item) for item in data]

    for group in groups:
        if group.id == id:
            return httpx.URL(f"{settings.SCHEDULE_LINK}{id}")

    return None
