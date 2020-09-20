import { RepoState } from '.';

export function startLoading(state: RepoState): void {
  state.loading = true;
}

export function stopLoading(state: RepoState): void {
  state.loading = false;
}

export function updateList(state: RepoState, repos: Repository[]): void {
  state.list.splice(0, state.list.length);
  state.list.push(...repos);
}

export function setCurrent(state: RepoState, repo: Repository): void {
  state.current = repo;
}

export function setError(state: RepoState, error?: Error): void {
  if (error) {
    state.error = error;
  } else {
    state.error = undefined;
  }
}

export function setSetting(state: RepoState, setting?: RepositorySetting): void {
  state.setting = setting;
}

export function setOwner(state: RepoState, owner: boolean): void {
  state.owner = owner;
}

export function setCommits(state: RepoState, commits: Commit[]): void {
  state.commits.splice(0, state.commits.length);
  state.commits.push(...commits);
}

export function setBranches(state: RepoState, branches: string[]): void {
  state.branches.splice(0, state.branches.length);
  state.branches.push(...branches);
}
