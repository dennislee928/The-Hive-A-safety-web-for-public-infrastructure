# ERH Safety System - Route 2 Flutter App

Flutter mobile application for the ERH Safety System Route 2 App, providing personalized safety guidance and crowd reporting capabilities.

## Features

- **Device Registration**: Automatic device registration with hashed device ID
- **Real-time Alerts**: CAP message display for safety alerts
- **Personalized Guidance**: Route recommendations to avoid unsafe zones
- **Crowd Reporting**: Submit safety reports with zone selection
- **Assistance Requests**: Request medical, security, or other assistance
- **Feedback Submission**: Submit feedback on guidance quality

## Requirements

- Flutter SDK >= 3.0.0
- Dart SDK >= 3.0.0
- Android SDK (for Android development)
- Xcode (for iOS development, macOS only)

## Installation

1. Install Flutter dependencies:
```bash
cd mobile_app
flutter pub get
```

2. For Android:
   - Configure Firebase (optional, for push notifications)
   - Update `android/app/build.gradle` with your package name

3. For iOS:
   - Run `pod install` in `ios/` directory
   - Configure Firebase (optional, for push notifications)

## Configuration

Update the API base URL in `lib/services/api_service.dart`:

```dart
static const String baseUrl = 'http://your-api-server:8080/api/v1';
```

For production, use HTTPS:
```dart
static const String baseUrl = 'https://your-api-server.com/api/v1';
```

## Running the App

### Development Mode
```bash
flutter run
```

### Release Mode
```bash
flutter run --release
```

### Build APK (Android)
```bash
flutter build apk --release
```

### Build IPA (iOS)
```bash
flutter build ios --release
```

## Project Structure

```
mobile_app/
├── lib/
│   ├── main.dart                 # App entry point
│   ├── models/                   # Data models
│   │   └── api_response.dart
│   ├── providers/                # State management (Provider)
│   │   ├── auth_provider.dart
│   │   ├── guidance_provider.dart
│   │   └── cap_provider.dart
│   ├── screens/                  # UI screens
│   │   ├── splash_screen.dart
│   │   ├── onboarding_screen.dart
│   │   ├── home_screen.dart
│   │   ├── report_screen.dart
│   │   ├── guidance_screen.dart
│   │   └── assistance_screen.dart
│   ├── services/                 # API and storage services
│   │   ├── api_service.dart
│   │   └── storage_service.dart
│   └── utils/                    # Utility functions
│       └── device_utils.dart
├── assets/                       # Images, icons
├── pubspec.yaml                  # Dependencies
└── README.md                     # This file
```

## Privacy & Security

- **Device ID Hashing**: Device IDs are hashed using SHA-256 before transmission
- **Secure Storage**: API keys stored in secure storage (encrypted)
- **Zone-based Location**: Only zone-level location data (Z1-Z4), no precise GPS
- **Data Minimization**: Only necessary data is collected and transmitted

## API Integration

The app integrates with the following API endpoints:

- `POST /api/v1/route2/devices/register` - Device registration
- `POST /api/v1/route2/devices/:device_id/push-token` - Push token registration
- `GET /api/v1/route2/guidance` - Get personalized guidance
- `POST /api/v1/reports` - Submit crowd report
- `POST /api/v1/route2/assistance` - Request assistance
- `POST /api/v1/route2/feedback` - Submit feedback
- `GET /api/v1/cap/zone/:zone_id` - Get CAP messages

## Testing

Run tests:
```bash
flutter test
```

## Build & Deployment

### Android

1. Generate keystore:
```bash
keytool -genkey -v -keystore ~/upload-keystore.jks -keyalg RSA -keysize 2048 -validity 10000 -alias upload
```

2. Configure signing in `android/app/build.gradle`

3. Build release APK:
```bash
flutter build apk --release
```

### iOS

1. Configure signing in Xcode
2. Build release:
```bash
flutter build ios --release
```

## Troubleshooting

### Common Issues

1. **API Connection Failed**
   - Check API server is running
   - Verify base URL in `api_service.dart`
   - Check network permissions in AndroidManifest.xml / Info.plist

2. **Build Errors**
   - Run `flutter clean`
   - Run `flutter pub get`
   - For iOS: `cd ios && pod install`

3. **Push Notifications Not Working**
   - Configure Firebase properly
   - Check device permissions
   - Verify push token registration

## License

See main project LICENSE file.

