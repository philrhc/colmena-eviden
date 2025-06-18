FROM grafana/grafana:latest

USER root

ENV GF_ENABLE_GZIP=true
ENV GF_USERS_DEFAULT_THEME=system
ENV GF_PANELS_DISABLE_SANITIZE_HTML=true
ENV GF_ANALYTICS_CHECK_FOR_UPDATES=false

#Install Plugin for SLA Panel
ENV GF_INSTALL_PLUGINS="yesoreyeram-infinity-datasource"
ENV GF_PLUGINS_ALLOW_LOADING_UNSIGNED="yesoreyeram-infinity-datasource"
ENV GF_FEATURE_TOGGLES_DATACONNECTIONSCONSOLE=true

# Explore
ENV GF_EXPLORE_ENABLED=true

# Alerting
ENV GF_ALERTING_ENABLED=false
ENV GF_UNIFIED_ALERTING_ENABLED=false

# Create necessary directories for Grafana data and dashboards
RUN mkdir -p /var/lib/grafana/data /usr/share/grafana/public/dashboards

# Copy your home.json dashboard to the appropriate location
COPY ./provisioning/dashboards/home.json /usr/share/grafana/public/dashboards/home.json

# Set the correct permissions for the dashboard file
RUN chmod 644 /usr/share/grafana/public/dashboards/home.json

COPY provisioning/dashboards.yaml /etc/grafana/provisioning/dashboards/dashboards.yaml

COPY grafana.ini /etc/grafana/grafana.ini
COPY grafana.db /var/lib/grafana/grafana.db

# Set UI elements 
COPY img/logo.svg /usr/share/grafana/public/img/grafana_icon.svg
COPY img/colmenas_favicon.png /usr/share/grafana/public/img/apple-touch-icon.png
COPY img/colmenas_favicon.png /usr/share/grafana/public/img/fav32.png
COPY img/colmenas_background.svg /usr/share/grafana/public/img/g8_login_dark.svg
COPY img/colmenas_background.svg /usr/share/grafana/public/img/g8_login_light.svg

RUN find /usr/share/grafana/public/build/ -name "*.js" \
    -exec sed -i 's|AppTitle="Grafana"|AppTitle="COLMENA Grafana"|g' {} \; \
    -exec sed -i 's|LoginTitle="Welcome to Grafana"|LoginTitle="Welcome to COLMENA"|g' {} \; \
    && find /usr/share/grafana/public/build/ -name "*.js" \
    -exec sed -i 's|.id==="enterprise"|.id==="notanenterprise"|g' {} \; \
    -exec sed -i 's|.id==="cloud"|.id==="notacloud"|g' {} \;

RUN sed -i "s|\[\[.NavTree\]\],|nav,|g; \
    s|window.grafanaBootData = {| \
    let nav = [[.NavTree]]; \
    const help = nav.find((element) => element.id === 'help'); \
    if (help) { help['subTitle'] = 'Business Suite 11.5.1'; help['children'] = [];} \
    const bookmarks = nav.find((element) => element.id === 'bookmarks'); \
    if (bookmarks) { nav = nav.filter((item) => item.id !== 'bookmarks'); } \
    const starred = nav.find((element) => element.id === 'starred'); \
    if (starred) { nav = nav.filter((item) => item.id !== 'starred'); } \
    window.grafanaBootData = {|g" \
/usr/share/grafana/public/views/index.html

## Add new route of menu for control panel
RUN sed -i "s|window.grafanaBootData = {| \
    nav.push({\"id\":\"my-settings\",\"text\":\"Control Panel\",\"icon\":\"brain\",\"url\":\"#\",\"sortWeight\":-1500}); \
    window.grafanaBootData = {|g" \
/usr/share/grafana/public/views/index.html



RUN find /usr/share/grafana/public/build/ -name *.js \
## Update Title
    -exec sed -i 's|AppTitle="Grafana"|AppTitle="COLMENA"|g' {} \; \
## Update Login Title
    -exec sed -i 's|LoginTitle="Welcome to COLMENA"|LoginTitle="COLMENA"|g' {} \; \
## Remove Documentation, Support, Community in the Footer
    -exec sed -i 's|\[{target:"_blank",id:"documentation".*grafana_footer"}\]|\[\]|g' {} \; \
## Remove Edition in the Footer
    -exec sed -i 's|({target:"_blank",id:"license",.*licenseUrl})|()|g' {} \; \
## Remove Version in the Footer
    -exec sed -i 's|({target:"_blank",id:"version",text:..versionString,url:D?"https://github.com/grafana/grafana/blob/main/CHANGELOG.md":void 0})|()|g' {} \; \
## Remove News icon
    -exec sed -i 's|(0,t.jsx)(...,{className:ge,onClick:.*,iconOnly:!0,icon:"rss","aria-label":"News"})|null|g' {} \; \
## Remove Old Dashboard page icon
    -exec sed -i 's|(0,t.jsx)(u.I,{tooltip:"Switch to old dashboard page",icon:"apps",onClick:()=>{s.Ny.partial({scenes:!1})}},"view-in-old-dashboard-button")|null|g' {} \; \
## Remove Open Source icon
    -exec sed -i 's|.push({target:"_blank",id:"version",text:`${..edition}${.}`,url:..licenseUrl,icon:"external-link-alt"})||g' {} \;

# Eliminar barra de búsqueda de la UI de Grafana
RUN find /usr/share/grafana/public/build/ -name "*.js" \
    -exec sed -i 's|.search-bar-container|/* .search-bar-container *|g' {} \; \
    -exec sed -i 's|search-bar__input|/* search-bar__input *|g' {} \; \
    -exec sed -i 's|search-bar__button|/* search-bar__button *|g' {} \; \
    -exec sed -i 's|search-bar__icon|/* search-bar__icon *|g' {} \; \
    -exec sed -i 's|.search-bar__clear|/* .search-bar__clear *|g' {} \;


## Remove datasources
RUN rm -rf /usr/share/grafana/public/app/plugins/datasource/elasticsearch /usr/share/grafana/public/build/elasticsearch* \
    /usr/share/grafana/public/app/plugins/datasource/graphite /usr/share/grafana/public/build/graphite* \
    /usr/share/grafana/public/app/plugins/datasource/opentsdb /usr/share/grafana/public/build/opentsdb* \
    /usr/share/grafana/public/app/plugins/datasource/influxdb /usr/share/grafana/public/build/influxdb* \
    /usr/share/grafana/public/app/plugins/datasource/mssql /usr/share/grafana/public/build/mssql* \
    /usr/share/grafana/public/app/plugins/datasource/mysql /usr/share/grafana/public/build/mysql* \
    /usr/share/grafana/public/app/plugins/datasource/tempo /usr/share/grafana/public/build/tempo* \
    /usr/share/grafana/public/app/plugins/datasource/jaeger /usr/share/grafana/public/build/jaeger* \
    /usr/share/grafana/public/app/plugins/datasource/zipkin /usr/share/grafana/public/build/zipkin* \
    /usr/share/grafana/public/app/plugins/datasource/azuremonitor /usr/share/grafana/public/build/azureMonitor* \
    /usr/share/grafana/public/app/plugins/datasource/cloudwatch /usr/share/grafana/public/build/cloudwatch* \
    /usr/share/grafana/public/app/plugins/datasource/cloud-monitoring /usr/share/grafana/public/build/cloudMonitoring* \
    /usr/share/grafana/public/app/plugins/datasource/parca /usr/share/grafana/public/build/parca* \
    /usr/share/grafana/public/app/plugins/datasource/phlare /usr/share/grafana/public/build/phlare* \
    /usr/share/grafana/public/app/plugins/datasource/grafana-pyroscope-datasource /usr/share/grafana/public/build/pyroscope*

## Remove Cloud and Enterprise categories
RUN find /usr/share/grafana/public/build/ -name *.js \
    -exec sed -i 's|.id==="enterprise"|.id==="notanenterprise"|g' {} \; \
    -exec sed -i 's|.id==="cloud"|.id==="notacloud"|g' {} \;
#Erase footer
RUN sed -i 's|\[{target:"_blank",id:"documentation".*grafana_footer"}\]|\[\]|g' /usr/share/grafana/public/views/index.html && \
    sed -i 's|({target:"_blank",id:"license",.*licenseUrl})|()|g' /usr/share/grafana/public/views/index.html && \
    sed -i 's|({target:"_blank",id:"version",text:..versionString,url:D?"https://github.com/grafana/grafana/blob/main/CHANGELOG.md":void 0})|()|g' /usr/share/grafana/public/views/index.html && \
    sed -i 's|(0,t.jsx)(...,{className:ge,onClick:.*,iconOnly:!0,icon:"rss","aria-label":"News"})|null|g' /usr/share/grafana/public/views/index.html && \
    sed -i 's|(0,t.jsx)(u.I,{tooltip:"Switch to old dashboard page",icon:"apps",onClick:()=>{s.Ny.partial({scenes:!1})}},"view-in-old-dashboard-button")|null|g' /usr/share/grafana/public/views/index.html && \
    sed -i 's|.push({target:"_blank",id:"version",text:`${..edition}${.}`,url:..licenseUrl,icon:"external-link-alt"})||g' /usr/share/grafana/public/views/index.html

#Install Infinity datasource for SLA Panel
RUN grafana-cli plugins install yesoreyeram-infinity-datasource
#RUN chown -R grafana:grafana /usr/share/grafana/public/app/plugins/yesoreyeram-infinity-datasource

# Disable Connections
#RUN sed -i 's|# feature2 = false|dataConnectionsConsole = false|g' /usr/share/grafana/conf/defaults.ini

EXPOSE 3000

CMD ["grafana-server", "--homepath=/usr/share/grafana", "--config=/etc/grafana/grafana.ini"]
