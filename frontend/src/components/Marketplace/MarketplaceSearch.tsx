import React, { useState, useEffect } from 'react';
import { TemplateMarketplaceItem, TemplateSearchRequest, TemplateSearchResult } from '../../types';
import './MarketplaceSearch.scss';

interface MarketplaceSearchProps {
  onSelectTemplate: (template: TemplateMarketplaceItem) => void;
  onInstallTemplate: (template: TemplateMarketplaceItem) => void;
}

const MarketplaceSearch: React.FC<MarketplaceSearchProps> = ({
  onSelectTemplate,
  onInstallTemplate,
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [searchType, setSearchType] = useState('');
  const [searchTags, setSearchTags] = useState('');
  const [sortBy, setSortBy] = useState('rating');
  const [sortOrder, setSortOrder] = useState('desc');
  const [searchResults, setSearchResults] = useState<TemplateSearchResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [categories, setCategories] = useState<string[]>([]);
  const [popularTags, setPopularTags] = useState<string[]>([]);
  const [currentPage, setCurrentPage] = useState(0);
  const [pageSize] = useState(12);

  useEffect(() => {
    loadCategories();
    loadPopularTags();
    performSearch(); // Initial search
  }, []);

  const loadCategories = async () => {
    try {
      const categoriesData = await window.go.app.App.GetMarketplaceCategories();
      setCategories(categoriesData || []);
    } catch (err) {
      console.error('Failed to load categories:', err);
    }
  };

  const loadPopularTags = async () => {
    try {
      const tagsData = await window.go.app.App.GetMarketplaceTags(20);
      setPopularTags(tagsData || []);
    } catch (err) {
      console.error('Failed to load tags:', err);
    }
  };

  const performSearch = async (page = 0) => {
    setLoading(true);
    setError(null);

    try {
      const searchRequest: TemplateSearchRequest = {
        query: searchQuery,
        type: searchType,
        tags: searchTags,
        author: '',
        minRating: '',
        sortBy: sortBy,
        sortOrder: sortOrder,
        limit: pageSize.toString(),
        offset: (page * pageSize).toString(),
      };

      const results = await window.go.app.App.SearchMarketplaceTemplates(searchRequest);
      setSearchResults(results);
      setCurrentPage(page);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Search failed');
      setSearchResults(null);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = () => {
    setCurrentPage(0);
    performSearch(0);
  };

  const handlePageChange = (newPage: number) => {
    performSearch(newPage);
  };

  const handleTagClick = (tag: string) => {
    const currentTags = searchTags.split(',').map(t => t.trim()).filter(t => t);
    if (!currentTags.includes(tag)) {
      setSearchTags(currentTags.concat(tag).join(', '));
    }
  };

  const handleInstall = async (template: TemplateMarketplaceItem) => {
    try {
      await onInstallTemplate(template);
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Installation failed');
    }
  };

  const renderTemplate = (template: TemplateMarketplaceItem) => (
    <div key={template.id} className="marketplace-template">
      <div className="marketplace-template__header">
        <h3 className="marketplace-template__title">{template.name}</h3>
        <div className="marketplace-template__badges">
          <span className={`marketplace-template__type marketplace-template__type--${template.type}`}>
            {template.type}
          </span>
          <span className="marketplace-template__version">v{template.version}</span>
        </div>
      </div>

      <div className="marketplace-template__content">
        <p className="marketplace-template__description">{template.description}</p>
        
        <div className="marketplace-template__meta">
          <div className="marketplace-template__author">
            <span>👤</span>
            <span>{template.author}</span>
          </div>
          <div className="marketplace-template__rating">
            <span>⭐</span>
            <span>{template.rating.toFixed(1)}</span>
          </div>
          <div className="marketplace-template__downloads">
            <span>📥</span>
            <span>{template.downloads.toLocaleString()}</span>
          </div>
        </div>

        {template.tags.length > 0 && (
          <div className="marketplace-template__tags">
            {template.tags.slice(0, 3).map(tag => (
              <span key={tag} className="marketplace-template__tag">
                {tag}
              </span>
            ))}
            {template.tags.length > 3 && (
              <span className="marketplace-template__tag-more">
                +{template.tags.length - 3} more
              </span>
            )}
          </div>
        )}
      </div>

      <div className="marketplace-template__actions">
        <button
          onClick={() => onSelectTemplate(template)}
          className="marketplace-template__action marketplace-template__action--view"
        >
          👁️ View Details
        </button>
        <button
          onClick={() => handleInstall(template)}
          className="marketplace-template__action marketplace-template__action--install"
        >
          📥 Install
        </button>
      </div>
    </div>
  );

  const renderPagination = () => {
    if (!searchResults || searchResults.total <= pageSize) return null;

    const totalPages = Math.ceil(searchResults.total / pageSize);
    const maxPages = 5;
    const startPage = Math.max(0, currentPage - Math.floor(maxPages / 2));
    const endPage = Math.min(totalPages, startPage + maxPages);

    return (
      <div className="marketplace-pagination">
        <button
          onClick={() => handlePageChange(currentPage - 1)}
          disabled={currentPage === 0}
          className="marketplace-pagination__button"
        >
          ← Previous
        </button>

        <div className="marketplace-pagination__pages">
          {Array.from({ length: endPage - startPage }, (_, i) => {
            const page = startPage + i;
            return (
              <button
                key={page}
                onClick={() => handlePageChange(page)}
                className={`marketplace-pagination__page ${
                  page === currentPage ? 'marketplace-pagination__page--active' : ''
                }`}
              >
                {page + 1}
              </button>
            );
          })}
        </div>

        <button
          onClick={() => handlePageChange(currentPage + 1)}
          disabled={currentPage >= totalPages - 1}
          className="marketplace-pagination__button"
        >
          Next →
        </button>
      </div>
    );
  };

  return (
    <div className="marketplace-search">
      <div className="marketplace-search__header">
        <h2>Template Marketplace</h2>
        <p>Discover and install templates from the community</p>
      </div>

      <div className="marketplace-search__filters">
        <div className="marketplace-search__filter-row">
          <div className="marketplace-search__search-input">
            <input
              type="text"
              placeholder="Search templates..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
            />
            <button onClick={handleSearch} className="marketplace-search__search-button">
              🔍
            </button>
          </div>

          <select
            value={searchType}
            onChange={(e) => setSearchType(e.target.value)}
            className="marketplace-search__select"
          >
            <option value="">All Types</option>
            <option value="default">Default</option>
            <option value="custom">Custom</option>
            <option value="plugin">Plugin</option>
          </select>

          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value)}
            className="marketplace-search__select"
          >
            <option value="rating">Rating</option>
            <option value="downloads">Downloads</option>
            <option value="name">Name</option>
            <option value="created">Created</option>
            <option value="updated">Updated</option>
          </select>

          <select
            value={sortOrder}
            onChange={(e) => setSortOrder(e.target.value)}
            className="marketplace-search__select"
          >
            <option value="desc">Descending</option>
            <option value="asc">Ascending</option>
          </select>
        </div>

        <div className="marketplace-search__filter-row">
          <input
            type="text"
            placeholder="Tags (comma-separated)"
            value={searchTags}
            onChange={(e) => setSearchTags(e.target.value)}
            className="marketplace-search__tags-input"
          />
        </div>

        {popularTags.length > 0 && (
          <div className="marketplace-search__popular-tags">
            <span>Popular tags:</span>
            {popularTags.slice(0, 10).map(tag => (
              <button
                key={tag}
                onClick={() => handleTagClick(tag)}
                className="marketplace-search__tag-button"
              >
                {tag}
              </button>
            ))}
          </div>
        )}
      </div>

      {error && (
        <div className="marketplace-search__error">
          <span>❌ {error}</span>
          <button onClick={() => setError(null)}>✕</button>
        </div>
      )}

      {loading && (
        <div className="marketplace-search__loading">
          <div className="marketplace-search__spinner">🔄</div>
          <span>Searching marketplace...</span>
        </div>
      )}

      {searchResults && !loading && (
        <div className="marketplace-search__results">
          <div className="marketplace-search__results-header">
            <h3>
              {searchResults.total} templates found
              {searchResults.searchTime && (
                <span className="marketplace-search__search-time">
                  ({searchResults.searchTime})
                </span>
              )}
            </h3>
          </div>

          {searchResults.items.length === 0 ? (
            <div className="marketplace-search__empty">
              <h4>No templates found</h4>
              <p>Try adjusting your search criteria or browse popular templates.</p>
            </div>
          ) : (
            <>
              <div className="marketplace-search__grid">
                {searchResults.items.map(renderTemplate)}
              </div>
              {renderPagination()}
            </>
          )}
        </div>
      )}
    </div>
  );
};

export default MarketplaceSearch;