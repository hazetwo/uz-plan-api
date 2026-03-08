import httpx

from app.config.settings import settings
from app.utils.id import get_url_by_id

MOCK_DATA = [
    {"group_id": "30560", "name": "INF1"},
    {"group_id": "30561", "name": "INF2"},
]


def test_get_url_by_id():
    assert get_url_by_id("30560", MOCK_DATA) == httpx.URL(
        f"{settings.SCHEDULE_LINK}30560"
    )


def test_get_url_by_id_none():
    assert get_url_by_id("34328", MOCK_DATA) is None


def test_get_url_by_id_invalid():
    assert get_url_by_id("invalid", MOCK_DATA) is None
