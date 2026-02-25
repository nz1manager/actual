import os
from django.core.wsgi import get_wsgi_application

# BU QATOR JUDA MUHIM: settings so'zi katta harf bilan yozilishi kerak
os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'core.settings')

application = get_wsgi_application()
