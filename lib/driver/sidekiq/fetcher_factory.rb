module Driver
  module Sidekiq
    class FetcherFactory
      def initialize(driver_queue_reader, sidekiq_queues)
        @driver_queue_reader = driver_queue_reader
        @sidekiq_queues = sidekiq_queues
      end

      def new(mgr, options)
        queue_fetch = QueueFetch.new(@driver_queue_reader)
        non_blocking_fetch = NonBlockingFetch.new(@sidekiq_queues)
        fetch = CompositeFetch.new(queue_fetch => 10, non_blocking_fetch => 1)
        Fetcher.new(mgr, fetch)
      end
    end
  end
end