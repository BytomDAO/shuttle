import pytest

from app import create_app


@pytest.fixture
def app():
    app = create_app(config_name='testing')
    return app
