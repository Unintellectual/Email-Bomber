import os
import base64
import json
import google.auth
from google.auth.transport.requests import Request
from google_auth_oauthlib.flow import InstalledAppFlow
from googleapiclient.discovery import build
from googleapiclient.errors import HttpError
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from dotenv import load_dotenv

load_dotenv()


CREDENTIALS_PATH = os.getenv('CREDENTIALS_PATH', 'credentials.json')
TOKEN_PATH = os.getenv('TOKEN_PATH', 'token.json')


SCOPES = ['https://www.googleapis.com/auth/gmail.send']

def get_credentials():
    """Gets valid user credentials from storage."""
    creds = None

    if os.path.exists(TOKEN_PATH):
        with open(TOKEN_PATH, 'r') as token:
            creds = json.load(token)

    if not creds or not creds.valid:
        if creds and creds.expired and creds.refresh_token:
            creds.refresh(Request())
        else:
            flow = InstalledAppFlow.from_client_secrets_file(
                CREDENTIALS_PATH, SCOPES)
            creds = flow.run_local_server(port=0)
        # Save the credentials for the next run
        with open(TOKEN_PATH, 'w') as token:
            json.dump(creds, token)
    
    return creds

def send_email(service, to, subject, body):
    """Send an email to the specified recipient."""
    try:
        message = MIMEMultipart()
        message['to'] = to
        message['from'] = 'Kalbo Kobu <comjoed00509@gmail.com>'
        message['subject'] = subject

        msg = MIMEText(body)
        message.attach(msg)

        raw_message = base64.urlsafe_b64encode(message.as_bytes()).decode()

        send_message = service.users().messages().send(
            userId='me', body={'raw': raw_message}).execute()

        print(f"Message Id: {send_message['id']}")
        return send_message
    except HttpError as error:
        print(f'An error occurred: {error}')

def main():
    creds = get_credentials()
    service = build('gmail', 'v1', credentials=creds)

    to = "comdamnsdunnns@gmail.com"
    subject = "Test Email"
    body = "Kobu Kalbo kalbo"

    for i in range(100):
        send_email(service, to, subject, body)
        print(f"Email #{i+1} sent successfully!")

if __name__ == '__main__':
    main()
