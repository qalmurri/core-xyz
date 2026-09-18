defmodule TrafficService.Repo do
  use Ecto.Repo,
    otp_app: :traffic_service,
    adapter: Ecto.Adapters.Postgres
end
