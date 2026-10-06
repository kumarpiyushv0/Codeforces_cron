# Codeforces to Google Calendar Sync Cron

A serverless cron job written in **Go** that automatically fetches upcoming contests from the public [Codeforces API](https://codeforces.com/api/contest.list) and synchronizes them into **Google Calendar**.

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

- **Automated Sync**: Fetches upcoming Codeforces contests twice a week (Mondays & Thursdays).
- **Parallel Division Merging**: Intelligently groups simultaneous rounds (e.g. Div. 1 + Div. 2) into a single clean event.
- **Idempotent & Deduplicated**: Uses deterministic event IDs so repeated cron runs update contest timings without creating duplicates.
- **Configurable Reminders**: Sets popup reminders (30 minutes and 10 minutes prior to contest start).
- **Multiple Hosting Options (100% Free)**:
  - **GitHub Actions Scheduled Cron** *(Zero server infrastructure, runs on scheduled runner)*
  - **Vercel Serverless Function** *(Serverless HTTP endpoint with `vercel.json` cron)*
  - **Local / Docker / CLI**

---

## 🚀 Setup Guide (For Running Your Own Instance)

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
6. Under **"Integrate calendar"**, copy your **Calendar ID** (e.g., `your_email@gmail.com` or `xyz@group.calendar.google.com`).

---

### Step 3: Deployment Options

#### Option A: GitHub Actions (Recommended — Free & Reliable)

1. Push this repository to GitHub (public or private).
2. In your GitHub repository, go to **Settings** > **Secrets and variables** > **Actions**.
3. Add the following repository secrets:
   - `GOOGLE_CREDENTIALS_JSON`: The entire raw JSON content of your downloaded service account key file.
   - `CALENDAR_ID`: Your Google Calendar ID (e.g. `your_email@gmail.com`).
   - `REMINDER_MINUTES`: (Optional, default `30`).
4. The workflow in [`.github/workflows/cron.yml`](.github/workflows/cron.yml) will automatically run twice a week (every Monday and Thursday at 02:25 UTC / 07:55 AM IST). You can also trigger it manually anytime under the **Actions** tab by clicking **Run workflow**.

---

#### Option B: Vercel Serverless

1. Install the [Vercel CLI](https://vercel.com/docs/cli) or import the repository in the Vercel dashboard.
2. In Vercel Project Settings > **Environment Variables**, add:
   - `GOOGLE_CREDENTIALS_JSON`: The raw JSON content of your service account key.
   - `CALENDAR_ID`: Your Google Calendar ID.
   - `CRON_SECRET`: *(Optional)* A secure token to protect the `/api/cron` endpoint.
3. Deploy:
   ```bash
   vercel deploy --prod
   ```
4. Vercel will trigger `/api/cron` according to the schedule in [`vercel.json`](vercel.json).

---

### Step 4: Local Testing

To test the sync script locally on your machine:

1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```
2. Set `GOOGLE_CREDENTIALS_JSON` (or `GOOGLE_APPLICATION_CREDENTIALS=./service_account.json`) and `CALENDAR_ID`.
3. Run the CLI tool:
   ```bash
   export $(cat .env | xargs) && go run cmd/cron/main.go
   ```

---

## 📁 Project Structure

```
.
├── .github/workflows/
│   └── cron.yml             # GitHub Actions scheduled cron workflow
├── api/
│   └── cron.go              # Vercel Serverless HTTP handler
├── cmd/cron/
│   └── main.go              # CLI entrypoint for local & GitHub Actions
├── internal/
│   ├── calendar/            # Google Calendar client and event upsert logic
│   ├── codeforces/          # Codeforces API client and contest filtering
│   └── syncer/              # Synchronization orchestrator with grouping logic
├── go.mod
├── go.sum
├── vercel.json              # Vercel Cron configuration
└── README.md
```
