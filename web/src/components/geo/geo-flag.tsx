import { useEffect, useRef, useState } from "react";

const FLAG_URI = "https://hatscripts.github.io/circle-flags/flags";

type GeoFlagProps = {
  value: string;
  immediate?: boolean;
  className?: string;
  size?: number;
};

export function GeoFlag({ value, immediate, className, size = 24 }: GeoFlagProps) {
  const container = useRef<HTMLSpanElement>(null);
  const [intersected, setIntersected] = useState(false);
  const visible = Boolean(immediate || intersected);

  let code = value.toLowerCase();
  if (!code || code === "other") {
    code = "xx";
  }
  const src = `${FLAG_URI}/${code}.svg`;

  useEffect(() => {
    if (immediate) return;
    const el = container.current;
    if (!el) return;
    const obs = new IntersectionObserver(
      ([e]) => {
        if (e?.isIntersecting) {
          setIntersected(true);
          obs.disconnect();
        }
      },
      { rootMargin: "80px" },
    );
    obs.observe(el);
    return () => obs.disconnect();
  }, [immediate]);

  return (
    <span ref={container} className={className}>
      {visible && (
        <img
          src={src}
          alt=""
          width={size}
          height={size}
          className="inline-block rounded-full"
          onError={(e) => {
            e.currentTarget.src = `${FLAG_URI}/xx.svg`;
          }}
        />
      )}
    </span>
  );
}
