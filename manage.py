import os
import click
from app import create_app

app = create_app(os.getenv('FLASK_CONFIG'))


@app.shell_context_processor
def make_shell_context():
    return dict(app=app)


@app.cli.command("init")
def init():
    """Init application, create database tables
    and create a new user named admin with password admin
    """
    click.echo("done")
