class Config(object):
    BYTOM_HOST = 'http://127.0.0.1:9888'
    BYTOM_WEBSOCKET = 'ws://localhost:9888/websocket-subscribe'


class DevelopmentConfig(Config):
    pass


class TestingConfig(Config):
    pass


class ProductionConfig(Config):
    pass


config = {
    'default': DevelopmentConfig,
    'testing': TestingConfig,
    'production': ProductionConfig
}
