from datetime import date, datetime, time

from app.api.core.exceptions import ParsingException


def parse_date(value: str) -> date:
    value = value.strip()
    try:
        return datetime.fromisoformat(value).date()
    except ValueError as exc:
        raise ParsingException(f"Date parsing error: {exc}")


def parse_time(value: str | None) -> time | None:
    if value is None:
        return None
    value = value.strip()
    try:
        return datetime.fromisoformat(value).time()
    except ValueError as exc:
        raise ParsingException(f"Time parsing error: {exc}")
