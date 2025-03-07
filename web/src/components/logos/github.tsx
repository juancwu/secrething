import GitHubLogoDark from "@/assets/github-mark.svg";

export function GitHubLogo({ className = "" }: { className?: string }) {
	return <img src={GitHubLogoDark} alt="GitHub" className={className} />;
}
