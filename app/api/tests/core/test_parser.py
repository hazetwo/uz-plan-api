from bs4 import BeautifulSoup

from app.api.core.parser import parse
from app.api.tests.utils.mock_html import MOCK_HTML


def test_parse():
    soup = BeautifulSoup(MOCK_HTML, "html.parser")

    parse(soup)
