FROM noonat/rbenv-nodenv
MAINTAINER Nathan Ostgard <noonat@phuce.com>

ENV RUBY_VERSION=2.1.6

RUN rbenv install $RUBY_VERSION && \
    CONFIGURE_OPTS="--disable-install-doc" rbenv global $RUBY_VERSION && \
    gem install bundler && \
    apt-get update -y && \
    apt-get install -y --no-install-recommends libmysqlclient-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY Gemfile /opt/src/Gemfile
COPY Gemfile.lock /opt/src/Gemfile.lock
WORKDIR /opt/src
RUN bundle install

COPY . /opt/src
EXPOSE 80
ENV VIRTUAL_HOST vcd.phuce.com

CMD ["bundle", "exec", "rackup", "-p", "80"]
