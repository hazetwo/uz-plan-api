FROM python:3.14-slim

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV UV_LINK_MODE=copy
ENV PATH="/workspace/.venv/bin:$PATH"

WORKDIR /workspace

RUN pip install --no-cache-dir --upgrade pip uv

COPY pyproject.toml uv.lock ./
COPY apps/api/pyproject.toml apps/api/uv.lock ./apps/api/
RUN uv sync --frozen --no-dev --package uz-plan-api

COPY apps/api ./apps/api

EXPOSE 8000

CMD ["uv", "run", "--frozen", "--package", "uz-plan-api", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000", "--app-dir", "apps/api"]
