get_env_var() {
  local env_file="$1"
  local var_name="$2"
  echo $(grep "^$var_name=" "$env_file" | cut -d '=' -f2-)
}


env_file=".env"

if [[ ! -f "$env_file" ]]; then
  echo ".env ファイルが見つかりません"
  exit 1
fi

GEMINI_PROJECT_ID=$(get_env_var "$env_file" "GEMINI_PROJECT_ID")
GEMINI_API_KEY=$(get_env_var "$env_file" "GEMINI_API_KEY")
LINE_BOT_CHANNEL_SECRET=$(get_env_var "$env_file" "LINE_BOT_CHANNEL_SECRET")
LINE_BOT_CHANNEL_TOKEN=$(get_env_var "$env_file" "LINE_BOT_CHANNEL_TOKEN")


env_vars=(
  "GEMINI_PROJECT_ID=$GEMINI_PROJECT_ID"
  "GEMINI_API_KEY=$GEMINI_API_KEY"
  "LINE_BOT_CHANNEL_SECRET=$LINE_BOT_CHANNEL_SECRET"
  "LINE_BOT_CHANNEL_TOKEN=$LINE_BOT_CHANNEL_TOKEN"
)


#gcloud arguments
app_name="go-http-function"
region="asia-northeast1"
function_name="LinbotCallback"
go_version="go123"
env_var_string=""
for var in "${env_vars[@]}"; do
  env_var_string+="$var,"
done


gcloud beta run deploy "$app_name" \
  --source . \
  --region "$region" \
  --function "$function_name" \
  --base-image "$go_version" \
  --allow-unauthenticated \
  --set-env-vars "$env_var_string" \
  --project "arumonogohan-app" 