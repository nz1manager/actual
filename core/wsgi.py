import os
import sys
from django.core.wsgi import get_wsgi_application

# Loyihaning asosiy papkasini Python yo'liga qo'shamiz
path = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
if path not in sys.path:
    sys.path.append(path)

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'core.settings')

try:
    application = get_wsgi_application()
except Exception as e:
    # Agar xato bo'lsa, logga chiqaradi
    print(f"WSGI Loading Error: {e}")
    raise e
