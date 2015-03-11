#!/usr/bin/python

import sys
sys.path += ['lib']
import flask
import re
import socket
import urllib
import requests
import pycountry
from geoip import geolite2 as geo

app = flask.Flask(__name__)

invalid = 'Invalid IP address/hostname'
mapuri = '//maps.google.com/maps?f=q&source=s_q&hl=en&geocode=&ie=UTF8&iwloc=A&output=embed&z={zoom}&q={query}'


@app.route('/')
@app.route('/<ip>')
def main(ip=None):
    errors = []
    if not ip:
        ip = flask.request.remote_addr
    if ip == '127.0.0.1':
        ip = requests.get('http://icanhazip.com').text.strip('\n')

    if not verify(ip):
        errors.append('Invalid domain or IP address')
        match = None
        mapurl = False
        data = False
    else:
        ip = verify(ip)
        match = geo.lookup(ip)
        data = results(match)
        if match:
            map_query = urllib.quote(str('%s,%s' % (data['subdivision'], data['country'])).strip(','))
            mapurl = mapuri.format(zoom=data['accuracy'], query=map_query)
        else:
            mapurl = False

    results(match)
    return flask.render_template('index.html', data=data, map=mapurl, query=ip, errors=errors)


@app.route('/api/<ip>')
@app.route('/api/<ip>/<methods>')
def api(ip, methods=None):
    if methods:
        methods = methods.lower().split(',')
    else:
        methods = []
    query = verify(ip)
    if not query:
        return flask.jsonify({'success': False, 'message': 'Invalid ip or domain'})
    match = geo.lookup(query)
    data = results(match)
    if not data:
        if methods:
            return ''
        return flask.jsonify({'success': False, 'message': 'No results found!'})
    else:
        if methods:
            output = []
            for method in methods:
                if method in data:
                    if data[method]:
                        if isinstance(data[method], list):
                            output.append(','.join(data[method]))
                        else:
                            output.append(data[method])
                    else:
                        output.append('N/A')
                else:
                    output.append('N/A')
            return '|'.join([str(i) for i in output])
        else:
            return flask.jsonify(data)


def verify(query):
    isip = True if re.match(r'^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$', query.strip()) else False
    isdomain = True if re.match(r'^[A-Za-z0-9.-]+$', query.strip()) else False
    if isdomain:
        try:
            query = socket.gethostbyname(query)
        except:
            return False
    if not isip and not isdomain:
        return False
    return query


def results(match):
    if not match:
        return False
    country = pycountry.countries.get(alpha2=match.country).name if match.country else ''

    if match.subdivisions and country:
        tmp = pycountry.subdivisions.get(code='%s-%s' % (match.country, list(match.subdivisions)[-1]))
        subdivisions = []
        subdivisions.append(tmp.name)
        if tmp.type == 'Country':
            accuracy = 3
        elif tmp.type == 'State':
            accuracy = 4
        elif tmp.type == 'District':
            accuracy = 5
        elif tmp.type == 'City':
            accuracy = 5
        else:
            accuracy = 5
        if tmp.parent:
            subdivisions.append(tmp.parent.name)
    else:
        subdivisions = ''
        accuracy = 1

    if match.location:
        latitude, longitude = match.location[0], match.location[1]
    else:
        latitude, longitude = '', ''

    timezone = match.timezone if match.timezone else ''

    data = {
        'country': country,
        'country_abbr': match.country,
        'continent': match.continent,
        'subdivision': subdivisions,
        'latitude': latitude,
        'longitude': longitude,
        'timezone': timezone,
        'accuracy': accuracy,
        'ip': match.ip,
        'hostname': socket.getfqdn(match.ip)
    }
    return data


@app.errorhandler(404)
@app.errorhandler(500)
def page_not_found(error):
    return flask.redirect('/')


@app.after_request
def add_header(response):
    """ Allow CORS based requests (if someone wanted to support putting it into javascript) """
    response.headers['Access-Control-Allow-Origin'] = '*'
    return response


if __name__ == '__main__':
    app.debug = True
    app.run(host='0.0.0.0', port=4001)
