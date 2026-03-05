from datetime import date

import pytest

from app.core.exceptions import ParsingException
from app.utils import (
    get_week_end,
    parse_date,
    parse_time,
)


def test_parse_date():
    date = "2023-10-10"
    result = parse_date(date)
    assert result is not None


def test_parse_date_invalid():
    bad_date = "2043-30-40"
    with pytest.raises(ParsingException):
        parse_date(bad_date)


def test_parse_time():
    time = "10:35"
    result = parse_time(time)
    assert result is not None


def test_parse_time_invalid():
    bad_time = "100:75"
    with pytest.raises(ParsingException):
        parse_time(bad_time)


def test_parse_time_none():
    null_time = None
    result = parse_time(null_time)
    assert result is None


def test_get_week_end():
    assert get_week_end(date(2026, 1, 1)) == date(2026, 1, 7)
