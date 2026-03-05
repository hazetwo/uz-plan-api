from bs4 import BeautifulSoup

from app.core.parser import parse_schedule
from app.tests.utils.mock_html import MOCK_HTML


def test_parse():
    soup = BeautifulSoup(MOCK_HTML, "html.parser")

    parse_schedule(soup)
