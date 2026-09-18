# lib/traffic_service_web/channels/user_socket.ex
defmodule TrafficServiceWeb.UserSocket do
  use Phoenix.Socket

  # Route channel "traffic:lobby" ke TrafficChannel
  channel "traffic:*", TrafficServiceWeb.TrafficChannel

  @impl true
  def connect(_params, socket, _connect_info) do
    {:ok, socket}
  end

  @impl true
  def id(_socket), do: nil
end
