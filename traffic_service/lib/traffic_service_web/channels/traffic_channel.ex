# lib/traffic_service_web/channels/traffic_channel.ex
defmodule TrafficServiceWeb.TrafficChannel do
  use Phoenix.Channel

  @impl true
  def join("traffic:lobby", _payload, socket) do
    {:ok, socket}
  end

  @impl true
  def handle_in("send_route", payload, socket) do
    route_data = %{
      id: "route-#{System.unique_integer([:positive])}",
      from: payload["from"],
      to: payload["to"],
      duration: payload["duration"] || 10,
      start_time: System.system_time(:millisecond)
    }

    broadcast!(socket, "new_route", route_data)
    {:reply, :ok, socket}
  end
end
