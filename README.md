# MultiTool GO - Android приложение на Go

MultiTool GO - это многофункциональное приложение для Android, написанное на языке Go с использованием Gio UI.

## Функции

- 🏠 **Главная страница** - навигация по приложению
- 🧮 **Калькулятор** - вычисление математических выражений
- ⛏️ **Анализ ошибок Minecraft** - два режима:
  - Обычный: анализ логов и выявление распространенных ошибок
  - ИИ: интеграция с Gemini/OpenRouter для умного анализа
- 📺 **YouTube Загрузчик** - скачивание видео в MP3/MP4 с выбором качества (1080p, 720p, 480p, 360p)
- ⚙️ **Настройки** - включение/отключение автообновления
- ℹ️ **О приложении** - информация о версии и функциях

## Структура проекта

```
multibotGO/
├── main.go           # Основной код приложения
├── go.mod            # Модуль Go
├── go.sum            # Зависимости
└── README.md         # Этот файл
```

## Требования

- Go 1.19 или выше
- Для сборки под Android:
  - Android SDK
  - Android NDK
  - gomobile

## Установка зависимостей

```bash
go mod tidy
```

## Сборка

### Для Linux/Desktop

```bash
go build -o multibotgo .
```

### Для Android (с использованием gomobile)

1. Установите gomobile:
```bash
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init
```

2. Соберите APK:
```bash
gomobile bind -target=android -androidapi=21 -o multibotgo.aar .
```

Или создайте полный APK:
```bash
gomobile apk -target=android -androidapi=21 -o multibotgo.apk .
```

## Способы создания APK без GitHub Actions

### Способ 1: Локальная сборка с gomobile

```bash
# Установка инструментов
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init

# Сборка APK
gomobile apk -target=android -androidapi=21 -o app.apk .
```

Требования:
- Установленный Android SDK
- Установленный Android NDK
- Переменные окружения ANDROID_HOME и ANDROID_NDK_HOME

### Способ 2: Использование Docker

Создайте `Dockerfile`:

```dockerfile
FROM golang:1.19

RUN apt-get update && apt-get install -y \
    openjdk-11-jdk \
    wget \
    unzip

ENV ANDROID_HOME=/opt/android-sdk
ENV ANDROID_NDK_HOME=$ANDROID_HOME/ndk

# Установка Android SDK command-line tools
RUN mkdir -p $ANDROID_HOME && \
    wget https://dl.google.com/android/repository/commandlinetools-linux-8512546_latest.zip -O /tmp/tools.zip && \
    unzip /tmp/tools.zip -d $ANDROID_HOME && \
    yes | $ANDROID_HOME/cmdline-tools/bin/sdkmanager --install "platform-tools" "platforms;android-30" "build-tools;30.0.3" "ndk;23.1.7779620"

RUN go install golang.org/x/mobile/cmd/gomobile@latest && \
    gomobile init

WORKDIR /app
COPY . .

RUN gomobile apk -target=android -androidapi=21 -o /output/app.apk .
```

Сборка:
```bash
docker build -t multibotgo-builder .
docker run --rm -v $(pwd):/output multibotgo-builder
```

### Способ 3: Использование Termux на Android устройстве

На Android устройстве с Termux:

```bash
pkg update && pkg upgrade
pkg install golang go-git
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init

git clone <ваш-репозиторий>
cd multibotGO
gomobile apk -target=android -androidapi=21 -o app.apk .
```

## Сборка с помощью Android Studio

### Вариант 1: Интеграция через AAR библиотеку

1. **Создайте AAR библиотеку:**
```bash
gomobile bind -target=android -androidapi=21 -o libmultibotgo.aar .
```

2. **В Android Studio:**
   - Создайте новый Android проект (Java/Kotlin)
   - Скопируйте `libmultibotgo.aar` в папку `app/libs/`
   - В `app/build.gradle` добавьте:
   ```gradle
   dependencies {
       implementation files('libs/libmultibotgo.aar')
   }
   ```
   - Используйте классы из Go библиотеки в вашем Java/Kotlin коде

### Вариант 2: Полная интеграция

1. **Создайте проект в Android Studio** с нативной активностью

2. **Экспортируйте Go код как библиотеку:**
```bash
gomobile bind -target=android -androidapi=21 -javapkg="org.multibotgo" -o multibotgo.aar .
```

3. **Добавьте зависимости в `app/build.gradle`:**
```gradle
android {
    // ...
}

dependencies {
    implementation fileTree(dir: 'libs', include: ['*.aar', '*.jar'])
    // Другие зависимости
}
```

4. **Инициализируйте Go в `MainActivity.java`:**
```java
import org.multibotgo.Multibotgo;

public class MainActivity extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        
        // Инициализация Go библиотеки
        Multibotgo.init(this);
        
        // Использование функций из Go
        String result = Multibotgo.someFunction();
    }
}
```

### Вариант 3: Использование CMake с Go как C библиотекой

Для продвинутой интеграции можно использовать CGo и CMake для компиляции Go кода как native library.

## Настройка переменных окружения

Для сборки Android убедитесь, что установлены:

```bash
export ANDROID_HOME=$HOME/Android/Sdk
export ANDROID_NDK_HOME=$ANDROID_HOME/ndk/<version>
export PATH=$PATH:$ANDROID_HOME/tools:$ANDROID_HOME/platform-tools
```

## Примечания

- Для работы YouTube загрузчика потребуется дополнительная интеграция с yt-dlp или аналогичной библиотекой
- ИИ анализ требует API ключи от Gemini или OpenRouter
- Автообновление требует реализации серверной части для проверки версий

## Лицензия

MIT
