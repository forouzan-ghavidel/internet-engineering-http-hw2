<div dir="rtl" align="right">

<h2 dir="rtl" align="right">سؤال ۱: محدود کردن دسترسی به سرویس در شبکه داخلی</h2>

<p dir="rtl" align="right">قرار دادن مستقیم یک سرویس قدیمی در دسترس کاربران خارجی می‌تواند خطرناک باشد؛ زیرا ممکن است این سرویس امکانات کافی برای احراز هویت و کنترل دسترسی نداشته باشد. برای محدود کردن دسترسی، یک <code dir="ltr">Reverse Proxy</code> در مقابل سرویس قرار می‌گیرد و تنها نقطهٔ ورود درخواست‌های خارجی می‌شود. در این حالت، پروکسی درخواست‌ها را دریافت می‌کند و پس از بررسی مجاز بودن آن‌ها، درخواست‌های مجاز را به سرویس اصلی یا <code dir="ltr">Backend</code> می‌فرستد. به این ترتیب، می‌توان کنترل دسترسی را بدون تغییر کد سرویس قدیمی انجام داد.</p>

<p dir="rtl" align="right">در معماری پیشنهادی، با استفاده از <code dir="ltr">Docker Compose</code>، پروکسی و سرویس <code dir="ltr">ghcr.io/mccutchen/go-httpbin</code> در دو کانتینر جداگانه اجرا می‌شوند. برای ارتباط میان آن‌ها، یک شبکهٔ داخلی با گزینهٔ <code dir="ltr">internal: true</code> تعریف می‌شود. سرویس اصلی فقط به این شبکه متصل است و روی پورت <code dir="ltr">8081</code> گوش می‌دهد. پروکسی نیز به همین شبکه متصل می‌شود تا بتواند با نام سرویس و پورت داخلی آن به سرویس اصلی دسترسی داشته باشد. همچنین، پروکسی به یک شبکهٔ جداگانه برای دریافت درخواست‌های خارجی متصل می‌شود.</p>

<p dir="rtl" align="right">پورت <code dir="ltr">8081</code> سرویس اصلی نباید با گزینهٔ <code dir="ltr">ports</code> روی میزبان منتشر شود؛ زیرا انتشار آن می‌تواند مسیری برای دسترسی مستقیم کاربران خارجی و دور زدن کنترل‌های پروکسی ایجاد کند. فقط پورت <code dir="ltr">80</code> پروکسی روی میزبان منتشر می‌شود. بنابراین، کاربر درخواست خود را به آدرس میزبان و پورت <code dir="ltr">80</code> ارسال می‌کند؛ پروکسی درخواست مجاز را از طریق شبکهٔ داخلی به سرویس اصلی می‌فرستد و پاسخ آن را به کاربر بازمی‌گرداند. در نتیجه، کاربران خارجی از طریق پروکسی با سرویس اصلی ارتباط دارند و پورت داخلی آن مستقیماً در اختیارشان قرار نمی‌گیرد.</p>

<h2 dir="rtl" align="right">سؤال ۲: فلسفه <code dir="ltr">TLS Termination</code></h2>

<p dir="rtl" align="right">یکی از راه‌های افزایش امنیت سرویس‌های قدیمی، استفاده از <code dir="ltr">TLS Termination</code> در کنار احراز هویت با <code dir="ltr">JWT</code> است. در این روش، به‌جای آنکه سرویس اصلی خودش ارتباط <code dir="ltr">HTTPS</code> را برقرار کند، این کار را <code dir="ltr">Reverse Proxy</code> انجام می‌دهد. پروکسی گواهی <code dir="ltr">TLS</code> را در اختیار دارد و ارتباط امن با کاربر برقرار می‌کند. سپس درخواست را رمزگشایی کرده و به سرویس داخلی می‌فرستد.</p>

<p dir="rtl" align="right">این روش باعث می‌شود بدون تغییر در کد سرویس قدیمی، ارتباط کاربر تا پروکسی رمزگذاری شود. مدیریت گواهی و تنظیمات امنیتی نیز در یک نقطه انجام می‌شود و نیازی نیست هر سرویس داخلی از <code dir="ltr">HTTPS</code> پشتیبانی کند. البته اگر پروکسی درخواست را از طریق <code dir="ltr">HTTP</code> به سرویس داخلی بفرستد، این بخش از مسیر رمزگذاری نشده است؛ به همین دلیل باید در یک شبکه داخلی و مورد اعتماد قرار داشته باشد.</p>

<p dir="rtl" align="right"><code dir="ltr">TLS</code> و <code dir="ltr">JWT</code> وظایف متفاوتی دارند و در کنار هم استفاده می‌شوند. <code dir="ltr">JWT</code> برای بررسی اعتبار توکن و کنترل دسترسی به سرویس به کار می‌رود، در حالی که <code dir="ltr">TLS</code> از محرمانگی و یکپارچگی داده‌ها در مسیر ارتباط کاربر با پروکسی محافظت می‌کند. به این ترتیب، هم می‌توان درخواست‌ها را پیش از رسیدن به سرویس بررسی کرد و هم از اطلاعات در زمان انتقال محافظت کرد.</p>

<h2 dir="rtl" align="right">سؤال ۳: پیاده‌سازی <code dir="ltr">Reverse Proxy</code> با احراز هویت <code dir="ltr">JWT</code></h2>

<p dir="rtl" align="right">در این بخش، یک <code dir="ltr">Reverse Proxy</code> با زبان <code dir="ltr">Go</code> پیاده‌سازی شده است که وظیفه دریافت درخواست‌های HTTP، بررسی توکن JWT و ارسال درخواست‌های معتبر به سرویس داخلی را بر عهده دارد. برای شبیه‌سازی سرویس قدیمی از <code dir="ltr">go-httpbin</code> استفاده شده است که در یک کانتینر جداگانه اجرا می‌شود.</p>

<p dir="rtl" align="right">پروکسی روی پورت <code dir="ltr">80</code> اجرا می‌شود و با استفاده از <code dir="ltr">httputil.NewSingleHostReverseProxy</code> درخواست‌ها را به آدرس داخلی <code dir="ltr">http://backend:8081</code> ارسال می‌کند. سرویس اصلی فقط به شبکه داخلی Docker متصل است و پورت آن روی میزبان منتشر نمی‌شود. در حالت پیش‌فرض، پورت پروکسی برای آزمایش محلی فقط روی آدرس <code dir="ltr">127.0.0.1</code> منتشر می‌شود. همچنین امکان تنظیم آدرس اتصال پروکسی برای دریافت درخواست از شبکه خارجی وجود دارد.</p>

<h3 dir="rtl" align="right">نحوه احراز هویت</h3>

<p dir="rtl" align="right">برای بررسی درخواست‌ها، یک Middleware در برنامه Go تعریف شده است. این بخش پیش از ارسال درخواست به سرویس اصلی، هدر <code dir="ltr">Authorization</code> را بررسی می‌کند و انتظار دارد توکن به شکل <code dir="ltr">Bearer</code> ارسال شده باشد.</p>

<p dir="rtl" align="right">اعتبارسنجی توکن با استفاده از کتابخانه <code dir="ltr">golang-jwt/jwt/v5</code> انجام می‌شود. در این پیاده‌سازی فقط الگوریتم <code dir="ltr">HS256</code> پذیرفته می‌شود و امضای توکن، زمان انقضا و محدودیت‌های زمانی آن بررسی می‌شوند. وجود <code dir="ltr">exp</code> الزامی است و <code dir="ltr">nbf</code> نیز در صورت وجود بررسی می‌شود. توکن از پارامترهای query دریافت نمی‌شود.</p>

<p dir="rtl" align="right">اگر توکن وجود نداشته باشد یا معتبر نباشد، پروکسی پاسخ <code dir="ltr">401 Unauthorized</code> برمی‌گرداند و درخواست به سرویس اصلی ارسال نمی‌شود. در صورت معتبر بودن توکن، هدر Authorization حذف شده و درخواست به سرویس داخلی فرستاده می‌شود.</p>

<p dir="rtl" align="right">کلید امضای JWT از متغیر محیطی <code dir="ltr">JWT_SECRET</code> خوانده می‌شود و در کد برنامه قرار ندارد.</p>

<h3 dir="rtl" align="right">نحوه اجرا</h3>

<p dir="rtl" align="right">پیش‌نیازها شامل Git برای دریافت مخزن، Docker با دسترسی به daemon، Docker Compose، Python 3 برای تولید کلید و توکن آزمایشی و curl برای ارسال درخواست هستند. نصب Go روی میزبان لازم نیست؛ Dockerfile برنامه را می‌سازد. فرمان‌ها در Bash و از ریشه مخزن اجرا می‌شوند. پورت <code dir="ltr">80</code> میزبان باید آزاد باشد.</p>

<div dir="ltr" align="left">

```bash
git clone https://github.com/forouzan-ghavidel/internet-engineering-http-hw2.git
cd internet-engineering-http-hw2
```

</div>

<p dir="rtl" align="right">پروژه با <code dir="ltr">Docker Compose</code> اجرا می‌شود. ابتدا فایل محلی <code dir="ltr">practical/.env</code> با کلیدی متشکل از حداقل ۳۲ بایت تصادفی امن و با کدگذاری Base64 ایجاد می‌شود. فرمان زیر فایل موجود را بازنویسی نمی‌کند و فایل جدید را با مجوز <code dir="ltr">0600</code> می‌سازد. این فایل در Git و Docker build context نادیده گرفته می‌شود؛ مقدار کلید نباید منتشر شود.</p>

<div dir="ltr" align="left">

```bash
python3 - <<'PYTHON'
import base64
import os
import secrets
from pathlib import Path

path = Path("practical/.env")
try:
    fd = os.open(path, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
except FileExistsError:
    print("Existing .env preserved.")
else:
    with os.fdopen(fd, "w") as file:
        file.write("JWT_SECRET=" + base64.b64encode(secrets.token_bytes(32)).decode() + "\n")
    print("Created .env with permissions 0600.")
PYTHON
```

</div>

<p dir="rtl" align="right">سپس پیکربندی بررسی، تصویر پروکسی ساخته و کانتینرها در حالت محلی راه‌اندازی می‌شوند. در محیط لینوکس این آزمایش، دانلود وابستگی در شبکه پیش‌فرض build با timeout مواجه شد؛ از این رو فرمان build زیر از شبکه میزبان استفاده می‌کند. این انتخاب فقط مربوط به زمان ساخت است و شبکه اجرای سرویس‌ها را تغییر نمی‌دهد. نام تصویر با نام پیش‌فرض پروژه Compose در پوشه practical هماهنگ است. فایل <code dir="ltr">.env</code> نسبت به محل <code dir="ltr">compose.yaml</code> پیدا می‌شود. اگر فایل موجود کلید معتبر نداشته باشد، پروکسی راه‌اندازی نمی‌شود؛ این فرمان آن را خودکار جایگزین نمی‌کند.</p>

<div dir="ltr" align="left">

```bash
unset PROXY_BIND_ADDRESS
docker compose -f practical/compose.yaml config --quiet
docker build --network host -t practical-proxy practical
docker compose -f practical/compose.yaml up -d --no-build
docker compose -f practical/compose.yaml ps
```

</div>

<h3 dir="rtl" align="right">تولید توکن و آزمایش درخواست‌ها</h3>

<p dir="rtl" align="right">فرمان زیر با کتابخانه استاندارد Python و همان کلید محلی، توکن HS256 دارای <code dir="ltr">exp</code> تولید می‌کند و آن را در متغیر پوسته <code dir="ltr">TOKEN</code> قرار می‌دهد. توکن ده دقیقه اعتبار دارد و پس از انقضا باید دوباره تولید شود. مقدار <code dir="ltr">TOKEN</code> یا secret را چاپ یا در مخزن ذخیره نکنید.</p>

<div dir="ltr" align="left">

```bash
TOKEN=$(python3 - <<'PYTHON'
import base64
import hashlib
import hmac
import json
import time
from pathlib import Path

value = Path("practical/.env").read_text().strip().split("=", 1)[1]
key = base64.b64decode(value, validate=True)

def encode(value):
    return base64.urlsafe_b64encode(value).rstrip(b"=").decode()

header = encode(json.dumps({"alg": "HS256", "typ": "JWT"}).encode())
payload = encode(json.dumps({"exp": int(time.time()) + 600}).encode())
message = header + "." + payload
signature = hmac.new(key, message.encode(), hashlib.sha256).digest()
print(message + "." + encode(signature))
PYTHON
)
```

</div>

<p dir="rtl" align="right">درخواست بدون توکن باید پاسخ <code dir="ltr">401</code> بگیرد. درخواست‌های دارای توکن معتبر باید پاسخ <code dir="ltr">200</code> دریافت کنند. گزینه <code dir="ltr">-i</code> کد وضعیت HTTP را نیز نمایش می‌دهد؛ پاسخ‌های GET و POST امکان بررسی query، بدنه و نبود هدر Authorization در سرویس اصلی را فراهم می‌کنند.</p>

<div dir="ltr" align="left">

```bash
curl --noproxy '*' -i http://127.0.0.1/get
curl --noproxy '*' -i -H "Authorization: Bearer $TOKEN" http://127.0.0.1/get
curl --noproxy '*' -i -H "Authorization: Bearer $TOKEN" 'http://127.0.0.1/get?course=http&item=one&item=two'
curl --noproxy '*' -i -H "Authorization: Bearer $TOKEN" -H 'Content-Type: text/plain' --data-binary 'hello from JWT proxy' http://127.0.0.1/post
curl --noproxy '*' -i -H "Authorization: Bearer $TOKEN" http://127.0.0.1/unknown-path
unset TOKEN
```

</div>

<h3 dir="rtl" align="right">حالت محلی و حالت خارجی</h3>

<p dir="rtl" align="right">نگاشت پورت پروکسی به صورت زیر تعریف شده است. در نبود متغیر <code dir="ltr">PROXY_BIND_ADDRESS</code>، مقدار پیش‌فرض <code dir="ltr">127.0.0.1</code> است. تنظیم <code dir="ltr">PROXY_BIND_ADDRESS=0.0.0.0</code> پورت <code dir="ltr">80</code> پروکسی را روی تمام رابط‌های IPv4 میزبان منتشر می‌کند. backend در هر دو حالت بدون پورت منتشرشده باقی می‌ماند؛ این محدودیت به معنی جلوگیری از تمام روش‌های دسترسی مدیر محلی Docker نیست.</p>

<div dir="ltr" align="left">

```yaml
ports:
  - "${PROXY_BIND_ADDRESS:-127.0.0.1}:80:80"
```

</div>

<p dir="rtl" align="right">حالت خارجی فقط باید در محیط آزمایشگاهی مورد اعتماد و با مجوز صریح اجرا شود. برای انتخاب آن هنگام راه‌اندازی، متغیر فوق پیش از فرمان Compose تنظیم می‌شود. در این مرحله حالت خارجی اجرا نشده و تنها اعتبار پیکربندی آن با فرمان زیر بررسی شده است؛ این فرمان چیزی راه‌اندازی نمی‌کند و secret را نمایش نمی‌دهد.</p>

<div dir="ltr" align="left">

```bash
PROXY_BIND_ADDRESS=0.0.0.0 docker compose -f practical/compose.yaml config --quiet
```

</div>

<p dir="rtl" align="right">این برنامه فعلا از HTTP استفاده می‌کند؛ بنابراین JWTهای واقعی نباید در شبکه غیرقابل اعتماد منتقل شوند. پیاده‌سازی TLS در این پروژه انجام نشده و مفهوم TLS Termination در پاسخ سؤال ۲ توضیح داده شده است.</p>

<h3 dir="rtl" align="right">بررسی عملکرد</h3>

<p dir="rtl" align="right">درخواست‌های بدون توکن یا با توکن نامعتبر با پاسخ <code dir="ltr">401 Unauthorized</code> رد می‌شوند. در مقابل، درخواست‌هایی که توکن معتبر دارند، به سرویس داخلی منتقل می‌شوند و پاسخ آن به کاربر بازگردانده می‌شود.</p>

<p dir="rtl" align="right">آزمایش‌های انجام‌شده شامل توکن منقضی‌شده، امضای نادرست، الگوریتم نامعتبر، درخواست GET و POST، پارامترهای query و مسیر ناموجود بوده‌اند. نتایج نشان دادند که درخواست‌های غیرمجاز پیش از رسیدن به سرویس اصلی رد می‌شوند و درخواست‌های معتبر به درستی پردازش می‌شوند.</p>

| آزمایش | کد وضعیت HTTP |
| --- | --- |
| بدون JWT | `401` |
| JWT بدقالب | `401` |
| امضای نادرست | `401` |
| توکن منقضی‌شده | `401` |
| بدون exp | `401` |
| الگوریتم غیرمجاز | `401` |
| GET با توکن معتبر | `200` |
| POST با توکن معتبر | `200` |
| مسیر ناموجود با توکن معتبر | `404` |

<p dir="rtl" align="right">برای اجرای تست‌های خودکار بدون نصب Go روی میزبان، می‌توان از تصویر build استفاده کرد. شبکه میزبان در این فرمان فقط برای دریافت وابستگی در کانتینر موقت تست است و شبکه سرویس‌ها را تغییر نمی‌دهد. تست‌ها شامل رد درخواست‌های نامعتبر، حذف Authorization و اعتبارسنجی کلید هستند.</p>

<div dir="ltr" align="left">

```bash
docker run --rm --network host -v "$PWD/practical:/app:ro" -w /app golang:1.27.3-alpine3.24 go test -v ./...
```

</div>

<h3 dir="rtl" align="right">مشاهده وضعیت و پایان اجرا</h3>

<p dir="rtl" align="right">فرمان اول وضعیت کانتینرها را نشان می‌دهد. فرمان دوم فقط کانتینرها و شبکه‌های همین پروژه Compose را متوقف و حذف می‌کند؛ فایل محلی <code dir="ltr">.env</code> حذف نمی‌شود. برای اجرای دوباره، فرمان راه‌اندازی بخش قبل تکرار می‌شود.</p>

<div dir="ltr" align="left">

```bash
docker compose -f practical/compose.yaml ps
docker compose -f practical/compose.yaml down
```

</div>

</div>
