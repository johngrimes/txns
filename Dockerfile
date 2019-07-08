FROM johngrimes/ruby

RUN apt-get update && \
    apt-get install -y git build-essential libpq-dev
RUN gem install bundler -v 1.13.7
COPY ./ /usr/src/txns/
WORKDIR /usr/src/txns
RUN bundle

ENTRYPOINT []
CMD ["bundle", "exec", "rackup"]
EXPOSE 9292
