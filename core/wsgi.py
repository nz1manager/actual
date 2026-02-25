import os
from django.core.wsgi import get_wsgi_application

os.environ.setdefault('django.settings_module', 'core.settings')

application = get_wsgi_application()
