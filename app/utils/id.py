import httpx

from app.config.settings import settings
from app.models import Group


def get_url_by_id(id: str, groups_data) -> httpx.URL | None:
    if not id.isdigit() or len(id) > settings.LONGEST_ID:
        return None

    groups = [Group.model_validate(item) for item in groups_data]

    for group in groups:
        if group.group_id == id:
            return httpx.URL(f"{settings.SCHEDULE_LINK}{id}")

    return None
