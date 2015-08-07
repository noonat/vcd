FROM noonat/ruby-node
MAINTAINER Nathan Ostgard <noonat@phuce.com>

RUN apt-get update -y && \
    apt-get install -y --no-install-recommends libmysqlclient-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY Gemfile /opt/src/Gemfile
COPY Gemfile.lock /opt/src/Gemfile.lock
WORKDIR /opt/src
RUN bundle install

COPY . /opt/src
EXPOSE 80

CMD ["bundle", "exec", "rackup", "-p", "80"]
