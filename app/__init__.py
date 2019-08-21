from flask import Flask
from app import api
from app.config import config
from flask_cors import CORS
# from app.extensions import db


def create_app(config_name=None):
    """Application factory, used to create application
    """
    app = Flask('app')

    CORS(app)
    configure_app(app, config_name)
    configure_extensions(app)
    register_blueprints(app)

    return app


def configure_app(app, config_name):
    """set configuration for application
    """
    # default configuration
    if config_name is None:
        config_name = 'default'
    app.config.from_object(config[config_name])


def configure_extensions(app):
    """configure flask extensions
    """
    # db.init_app(app)


def register_blueprints(app):
    """register all blueprints for application
    """
    app.register_blueprint(api.blueprint)
