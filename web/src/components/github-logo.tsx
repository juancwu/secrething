import { useTheme } from "@/providers/theme-provider";
import GitHubLogoWhite from "@/assets/github-mark-white.png";
import GitHubLogoDark from "@/assets/github-mark.svg";
import { useEffect, useState } from "react";

export function GitHubLogo({ className = "" }: { className?: string }) {
  const { theme } = useTheme();
  const [effectiveTheme, setEffectiveTheme] = useState(theme);
  
  useEffect(() => {
    // Handle system theme
    if (theme === "system") {
      const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
      const updateTheme = (e: MediaQueryListEvent | MediaQueryList) => {
        setEffectiveTheme(e.matches ? "dark" : "light");
      };
      
      // Set initial value
      updateTheme(mediaQuery);
      
      // Listen for changes
      mediaQuery.addEventListener("change", updateTheme);
      return () => mediaQuery.removeEventListener("change", updateTheme);
    } else {
      setEffectiveTheme(theme);
    }
  }, [theme]);
  
  // Use white logo for dark theme, dark logo for light theme
  const logoSrc = effectiveTheme === "dark" ? GitHubLogoWhite : GitHubLogoDark;
  
  return (
    <img src={logoSrc} alt="GitHub" className={className} />
  );
}