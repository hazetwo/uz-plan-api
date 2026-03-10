import httpx

from app.config.settings import settings
from app.core.handlers.exceptions import UrlNotFoundException
from app.models import Group


def get_url_by_id(id: str, groups_data) -> httpx.URL:
    if not id.isdigit() or len(id) > settings.LONGEST_ID:
        raise UrlNotFoundException("Invalid ID.")

    groups = [Group.model_validate(item) for item in groups_data]

    for group in groups:
        if group.group_id == id:
            return httpx.URL(f"{settings.SCHEDULE_LINK}{id}")

    raise UrlNotFoundException("The provided ID does not map to a valid URL.")
