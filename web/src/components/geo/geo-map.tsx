import L from "leaflet";
import { useEffect, useMemo } from "react";
import { Circle, MapContainer, Marker, ScaleControl, TileLayer, useMap } from "react-leaflet";
import type { GeoResult } from "@/api/types.gen";

import "leaflet/dist/leaflet.css";

import markerIcon2x from "leaflet/dist/images/marker-icon-2x.png";
import markerShadow from "leaflet/dist/images/marker-shadow.png";

const TILE_URL = "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png";
const ATTRIBUTION = '&copy; <a href="https://www.openstreetmap.org/copyright">OSM</a>';

function FixResize() {
  const map = useMap();
  useEffect(() => {
    const t = setTimeout(() => map.invalidateSize(), 400);
    return () => clearTimeout(t);
  }, [map]);
  return null;
}

export function GeoMap({ value }: { value: GeoResult }) {
  const { zoom, showCircle } = useMemo(() => {
    let z = 5;
    const acc = value.accuracy_radius_km;
    if (acc > 0) {
      if (acc <= 25) z = 9;
      else if (acc <= 50) z = 8;
      else if (acc <= 200) z = 6;
      else if (acc <= 500) z = 5;
      else if (acc <= 1000) z = 3;
      else z = 2;
    }
    return { zoom: z, showCircle: acc > 0 };
  }, [value.accuracy_radius_km]);

  useEffect(() => {
    const DefaultIcon = L.Icon.Default.prototype as unknown as {
      _getIconUrl?: string;
    };
    delete DefaultIcon._getIconUrl;
    L.Icon.Default.mergeOptions({
      iconUrl: markerIcon2x,
      iconRetinaUrl: markerIcon2x,
      shadowUrl: markerShadow,
    });
  }, []);

  const center: [number, number] = [value.latitude, value.longitude];

  return (
    <MapContainer
      center={center}
      zoom={zoom}
      className="leaflet-dark-tiles z-0 h-[200px] w-full rounded-md contain-[layout]"
      scrollWheelZoom
      preferCanvas
    >
      <TileLayer attribution={ATTRIBUTION} url={TILE_URL} maxZoom={18} />
      <ScaleControl />
      <FixResize />
      {showCircle ? (
        <Circle
          center={center}
          radius={value.accuracy_radius_km * 1000}
          pathOptions={{
            weight: 1.5,
            color: "red",
            fillColor: "red",
            fillOpacity: 0.05,
          }}
        />
      ) : (
        <Marker position={center} />
      )}
    </MapContainer>
  );
}
