# Codeforces to Google Calendar Sync Cron

A scheduled automation written in **Go** that automatically fetches upcoming contests from the public [Codeforces API](https://codeforces.com/api/contest.list) and synchronizes them into **Google Calendar** using **GitHub Actions**.

---

## 📅 Subscribe to the Live Calendar (1-Click)

Never miss a Codeforces round! You can subscribe to this live-updating calendar on any device:

[![Add to Google Calendar](https://img.shields.io/badge/Google_Calendar-Add_Calendar-4285F4?style=for-the-badge&logo=google-calendar&logoColor=white)](https://calendar.google.com/calendar/r?cid=dc8e239bbfa6913e2fe9823a850cf76ee9f7781cbd81f9150bf5c99df91ed7df@group.calendar.google.com)

* **Google Calendar 1-Click Subscription**: [Click here to add to Google Calendar](https://calendar.google.com/calendar/r?cid=dc8e239bbfa6913e2fe9823a850cf76ee9f7781cbd81f9150bf5c99df91ed7df@group.calendar.google.com)
* **iCal Format (Apple Calendar / Outlook)**:
  ```text
  https://calendar.google.com/calendar/ical/dc8e239bbfa6913e2fe9823a850cf76ee9f7781cbd81f9150bf5c99df91ed7df%40group.calendar.google.com/public/basic.ics
  ```

---

## ✨ Features

- **Automated GitHub Actions Cron**: Runs automatically twice a week (every Monday & Thursday at 02:25 UTC / 07:55 AM IST).
- **Parallel Division Merging**: Intelligently groups simultaneous rounds (e.g. Div. 1 + Div. 2) into a single clean event.
- **Idempotent & Deduplicated**: Uses deterministic event IDs so repeated runs update contest schedules without creating duplicates.
- **Configurable Reminders**: Sets popup reminders (30 minutes and 10 minutes prior to contest start).
- **Zero Infrastructure Cost**: 100% free using GitHub Actions.

---

## 🚀 Setup & Deployment Guide

### Step 1: Google Cloud Service Account Setup

1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Create a new project (e.g., `Codeforces-Calendar-Sync`) or select an existing one.
3. Enable the **Google Calendar API**:
   - Navigate to **APIs & Services** > **Library**.
   - Search for **Google Calendar API** and click **Enable**.
4. Create a **Service Account**:
   - Go to **APIs & Services** > **Credentials**.
   - Click **Create Credentials** > **Service Account**.
   - Give it a name (e.g., `cf-sync-cron`) and click **Done**.
5. Generate a **JSON Key**:
   - Click on the created service account email.
   - Go to the **Keys** tab > **Add Key** > **Create new key** > Select **JSON** > **Create**.
   - Save the downloaded JSON file.

---

### Step 2: Share Google Calendar with the Service Account

1. Open [Google Calendar](https://calendar.google.com/).
2. Under "My calendars", hover over your target calendar and click **⋮ (Settings and sharing)**.
3. Scroll down to **"Share with specific people or groups"** and click **Add people and groups**.
4. Paste the **Service Account Email** (e.g., `cf-sync-cron@your-project.iam.gserviceaccount.com`).
5. Set permission to **"Make changes to events"** and click **Send**.
6. Under **"Integrate calendar"**, copy your **Calendar ID** (e.g., `your_calendar_id@group.calendar.google.com`).

---

### Step 3: Configure GitHub Actions Secrets

1. Push this repository to your GitHub account.
2. In your GitHub repository, go to **Settings** > **Secrets and variables** > **Actions**.
3. Add the following repository secrets:
   - `GOOGLE_CREDENTIALS_JSON`: The entire raw JSON content of your downloaded service account key file.
   - `CALENDAR_ID`: Your Google Calendar ID.
   - `REMINDER_MINUTES`: (Optional, default `30`).
4. The workflow in [`.github/workflows/cron.yml`](.github/workflows/cron.yml) will automatically run on schedule. You can also trigger it manually anytime under the **Actions** tab by clicking **Run workflow**.

---

### Step 4: Local Testing (Optional)

To run the sync script locally on your machine:

1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```
2. Set `GOOGLE_CREDENTIALS_JSON` and `CALENDAR_ID`.
3. Run the Go CLI:
   ```bash
   export $(cat .env | xargs) && go run cmd/cron/main.go
   ```

---

## 📁 Project Structure

```
.
├── .github/workflows/
│   └── cron.yml             # Scheduled GitHub Actions workflow (twice a week)
├── cmd/cron/
│   └── main.go              # CLI entrypoint for GitHub Actions & local runs
├── internal/
│   ├── calendar/            # Google Calendar client and event upsert logic
│   ├── codeforces/          # Codeforces API client and contest filtering
│   └── syncer/              # Synchronization orchestrator with division grouping logic
├── go.mod
├── go.sum
├── .env.example             # Environment variable template
└── README.md                # Documentation and subscription links
```
