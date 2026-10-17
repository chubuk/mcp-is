package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/chubuk/mcp-is/internal/core"
	"github.com/chubuk/mcp-is/internal/integration/insales"
)

type Service struct {
	manager *core.AccountManager
}

func NewService(mgr *core.AccountManager) *Service {
	return &Service{manager: mgr}
}

func (s *Service) checkReadOnly(ctx context.Context) error {
	if readOnly, ok := ctx.Value(core.ShopReadOnlyKey).(bool); ok && readOnly {
		return fmt.Errorf("shop is in read-only mode")
	}
	return nil
}

func (s *Service) GetStoreDetails(ctx context.Context) (*insales.StoreInfo, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetStoreInfo(ctx)
}

// === ТОВАРЫ ===

func (s *Service) GetProducts(ctx context.Context, page, perPage int, filters map[string]any) ([]insales.Product, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListProducts(ctx, page, perPage, filters)
}

func (s *Service) GetProduct(ctx context.Context, productID string) (*insales.Product, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetProduct(ctx, productID)
}

func (s *Service) CreateProduct(ctx context.Context, title, price string, args map[string]any) (*insales.Product, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateProduct(ctx, title, price, args)
}

func (s *Service) UpdateProduct(ctx context.Context, productID string, args map[string]any) (*insales.Product, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateProduct(ctx, productID, args)
}

func (s *Service) DeleteProduct(ctx context.Context, productID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteProduct(ctx, productID)
}

// === ЗАКАЗЫ ===

func (s *Service) GetOrders(ctx context.Context, page, perPage int, filters map[string]any) ([]insales.Order, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListOrders(ctx, page, perPage, filters)
}

func (s *Service) GetOrderDetails(ctx context.Context, orderID string) (*insales.Order, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetOrder(ctx, orderID)
}

func (s *Service) UpdateOrder(ctx context.Context, orderID string, args map[string]any) (*insales.Order, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateOrder(ctx, orderID, args)
}

func (s *Service) CreateOrder(ctx context.Context, clientID int, args map[string]any) (*insales.Order, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateOrder(ctx, clientID, args)
}

// === КЛИЕНТЫ ===

func (s *Service) GetClients(ctx context.Context, page, perPage int) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListClients(ctx, page, perPage)
}

func (s *Service) GetClient(ctx context.Context, clientID string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetClient(ctx, clientID)
}

func (s *Service) CreateClient(ctx context.Context, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateClient(ctx, args)
}

func (s *Service) UpdateClient(ctx context.Context, clientID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateClient(ctx, clientID, args)
}

func (s *Service) DeleteClient(ctx context.Context, clientID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteClient(ctx, clientID)
}

// === КАТЕГОРИИ ===

func (s *Service) GetCategories(ctx context.Context, page, perPage int) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListCategories(ctx, page, perPage)
}

func (s *Service) GetCategory(ctx context.Context, categoryID string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetCategory(ctx, categoryID)
}

func (s *Service) CreateCategory(ctx context.Context, title string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateCategory(ctx, title, args)
}

func (s *Service) UpdateCategory(ctx context.Context, categoryID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateCategory(ctx, categoryID, args)
}

func (s *Service) DeleteCategory(ctx context.Context, categoryID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteCategory(ctx, categoryID)
}

// === КОЛЛЕКЦИИ ===

func (s *Service) GetCollections(ctx context.Context, page, perPage int) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListCollections(ctx, page, perPage)
}

func (s *Service) GetCollection(ctx context.Context, collectionID string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetCollection(ctx, collectionID)
}

func (s *Service) CreateCollection(ctx context.Context, title string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateCollection(ctx, title)
}

func (s *Service) UpdateCollection(ctx context.Context, collectionID, title string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateCollection(ctx, collectionID, title)
}

func (s *Service) DeleteCollection(ctx context.Context, collectionID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteCollection(ctx, collectionID)
}

// === ИЗОБРАЖЕНИЯ ===

func (s *Service) ListImages(ctx context.Context, productID string) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListImages(ctx, productID)
}

func (s *Service) AddImage(ctx context.Context, productID, url string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.AddImage(ctx, productID, url, args)
}

func (s *Service) DeleteImage(ctx context.Context, productID, imageID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteImage(ctx, productID, imageID)
}

// === МЕТАПОЛЯ ===

func (s *Service) ListMetaFields(ctx context.Context, args map[string]any) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListMetaFields(ctx, args)
}

func (s *Service) CreateMetaField(ctx context.Context, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateMetaField(ctx, args)
}

func (s *Service) UpdateMetaField(ctx context.Context, metaFieldID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateMetaField(ctx, metaFieldID, args)
}

func (s *Service) DeleteMetaField(ctx context.Context, metaFieldID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteMetaField(ctx, metaFieldID)
}

// === ВЕБХУКИ ===

func (s *Service) ListWebhooks(ctx context.Context) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListWebhooks(ctx)
}

func (s *Service) CreateWebhook(ctx context.Context, event, url string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateWebhook(ctx, event, url, args)
}

func (s *Service) UpdateWebhook(ctx context.Context, webhookID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateWebhook(ctx, webhookID, args)
}

func (s *Service) DeleteWebhook(ctx context.Context, webhookID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteWebhook(ctx, webhookID)
}

// === СИСТЕМНЫЕ ===

func (s *Service) ListCurrencies(ctx context.Context) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListCurrencies(ctx)
}

func (s *Service) ListDeliveryVariants(ctx context.Context) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListDeliveryVariants(ctx)
}

func (s *Service) ListPaymentGateways(ctx context.Context) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListPaymentGateways(ctx)
}

// === СТРАНИЦЫ ===

func (s *Service) ListPages(ctx context.Context, page, perPage int) ([]insales.Page, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListPages(ctx, page, perPage)
}

func (s *Service) GetPage(ctx context.Context, pageID string) (*insales.Page, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetPage(ctx, pageID)
}

func (s *Service) CreatePage(ctx context.Context, title, content, handle string) (*insales.Page, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreatePage(ctx, title, content, handle)
}

func (s *Service) UpdatePage(ctx context.Context, pageID string, args map[string]any) (*insales.Page, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdatePage(ctx, pageID, args)
}

func (s *Service) DeletePage(ctx context.Context, pageID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeletePage(ctx, pageID)
}

// === АКСЕССУАРЫ ===

func (s *Service) ListAccessories(ctx context.Context, page, perPage int) ([]insales.Accessory, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListAccessories(ctx, page, perPage)
}

func (s *Service) GetAccessory(ctx context.Context, accessoryID string) (*insales.Accessory, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetAccessory(ctx, accessoryID)
}

func (s *Service) CreateAccessory(ctx context.Context, title string, args map[string]any) (*insales.Accessory, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateAccessory(ctx, title, args)
}

func (s *Service) UpdateAccessory(ctx context.Context, accessoryID string, args map[string]any) (*insales.Accessory, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateAccessory(ctx, accessoryID, args)
}

func (s *Service) DeleteAccessory(ctx context.Context, accessoryID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteAccessory(ctx, accessoryID)
}

// === ВАРИАНТЫ ===

func (s *Service) ListVariants(ctx context.Context, productID string) ([]insales.Variant, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListVariants(ctx, productID)
}

func (s *Service) GetVariant(ctx context.Context, productID, variantID string) (*insales.Variant, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetVariant(ctx, productID, variantID)
}

func (s *Service) UpdateVariant(ctx context.Context, productID, variantID string, args map[string]any) (*insales.Variant, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateVariant(ctx, productID, variantID, args)
}

func (s *Service) DeleteVariant(ctx context.Context, productID, variantID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteVariant(ctx, productID, variantID)
}

// === ЗНАЧЕНИЯ СВОЙСТВ ===

func (s *Service) ListPropertyValues(ctx context.Context, accessoryID string) ([]insales.PropertyValue, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListPropertyValues(ctx, accessoryID)
}

func (s *Service) CreatePropertyValue(ctx context.Context, accessoryID, title string) (*insales.PropertyValue, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreatePropertyValue(ctx, accessoryID, title)
}

func (s *Service) UpdatePropertyValue(ctx context.Context, accessoryID, propertyValueID string, args map[string]any) (*insales.PropertyValue, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdatePropertyValue(ctx, accessoryID, propertyValueID, args)
}

func (s *Service) DeletePropertyValue(ctx context.Context, accessoryID, propertyValueID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeletePropertyValue(ctx, accessoryID, propertyValueID)
}

// === МЕНЮ ===

func (s *Service) ListMenus(ctx context.Context) ([]insales.Menu, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListMenus(ctx)
}

func (s *Service) GetMenu(ctx context.Context, menuID string) (*insales.Menu, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetMenu(ctx, menuID)
}

func (s *Service) CreateMenu(ctx context.Context, title, handle string) (*insales.Menu, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateMenu(ctx, title, handle)
}

func (s *Service) UpdateMenu(ctx context.Context, menuID string, args map[string]any) (*insales.Menu, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateMenu(ctx, menuID, args)
}

func (s *Service) DeleteMenu(ctx context.Context, menuID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteMenu(ctx, menuID)
}

// === РЕДИРЕКТЫ ===

func (s *Service) ListRedirects(ctx context.Context) ([]insales.Redirect, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListRedirects(ctx)
}

func (s *Service) CreateRedirect(ctx context.Context, url, target string) (*insales.Redirect, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateRedirect(ctx, url, target)
}

func (s *Service) DeleteRedirect(ctx context.Context, redirectID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteRedirect(ctx, redirectID)
}

// === БЛОГИ ===

func (s *Service) ListBlogs(ctx context.Context) ([]insales.Blog, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListBlogs(ctx)
}

func (s *Service) GetBlog(ctx context.Context, blogID string) (*insales.Blog, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetBlog(ctx, blogID)
}

func (s *Service) CreateBlog(ctx context.Context, title, handle string) (*insales.Blog, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateBlog(ctx, title, handle)
}

func (s *Service) UpdateBlog(ctx context.Context, blogID string, args map[string]any) (*insales.Blog, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateBlog(ctx, blogID, args)
}

func (s *Service) DeleteBlog(ctx context.Context, blogID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteBlog(ctx, blogID)
}

func (s *Service) ListBlogPosts(ctx context.Context, blogID string) ([]insales.BlogPost, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListBlogPosts(ctx, blogID)
}

func (s *Service) GetBlogPost(ctx context.Context, blogID, postID string) (*insales.BlogPost, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetBlogPost(ctx, blogID, postID)
}

func (s *Service) CreateBlogPost(ctx context.Context, blogID, title, bodyHTML string) (*insales.BlogPost, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateBlogPost(ctx, blogID, title, bodyHTML)
}

func (s *Service) UpdateBlogPost(ctx context.Context, blogID, postID string, args map[string]any) (*insales.BlogPost, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateBlogPost(ctx, blogID, postID, args)
}

func (s *Service) DeleteBlogPost(ctx context.Context, blogID, postID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteBlogPost(ctx, blogID, postID)
}

// === DISCOUNT CODES ===

func (s *Service) ListDiscountCodes(ctx context.Context, page, perPage int) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListDiscountCodes(ctx, page, perPage)
}

func (s *Service) GetDiscountCode(ctx context.Context, codeID string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetDiscountCode(ctx, codeID)
}

func (s *Service) CreateDiscountCode(ctx context.Context, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateDiscountCode(ctx, args)
}

func (s *Service) UpdateDiscountCode(ctx context.Context, codeID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateDiscountCode(ctx, codeID, args)
}

func (s *Service) DeleteDiscountCode(ctx context.Context, codeID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteDiscountCode(ctx, codeID)
}

// === COLLECTS ===

func (s *Service) ListCollects(ctx context.Context, page, perPage int) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListCollects(ctx, page, perPage)
}

func (s *Service) GetCollectsByProduct(ctx context.Context, productID string) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetCollectsByProduct(ctx, productID)
}

func (s *Service) CreateCollect(ctx context.Context, collectionID, productID int) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateCollect(ctx, collectionID, productID)
}

func (s *Service) UpdateCollect(ctx context.Context, collectID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateCollect(ctx, collectID, args)
}

func (s *Service) DeleteCollect(ctx context.Context, collectID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteCollect(ctx, collectID)
}

// === PROPERTIES ===

func (s *Service) ListProperties(ctx context.Context) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListProperties(ctx)
}

func (s *Service) GetProperty(ctx context.Context, propertyID string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetProperty(ctx, propertyID)
}

func (s *Service) CreateProperty(ctx context.Context, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateProperty(ctx, args)
}

func (s *Service) UpdateProperty(ctx context.Context, propertyID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateProperty(ctx, propertyID, args)
}

func (s *Service) DeleteProperty(ctx context.Context, propertyID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteProperty(ctx, propertyID)
}

// === CHARACTERISTICS ===

func (s *Service) ListCharacteristics(ctx context.Context, propertyID string) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListCharacteristics(ctx, propertyID)
}

func (s *Service) GetCharacteristic(ctx context.Context, propertyID, characteristicID string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetCharacteristic(ctx, propertyID, characteristicID)
}

func (s *Service) CreateCharacteristic(ctx context.Context, propertyID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateCharacteristic(ctx, propertyID, args)
}

func (s *Service) UpdateCharacteristic(ctx context.Context, propertyID, characteristicID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateCharacteristic(ctx, propertyID, characteristicID, args)
}

func (s *Service) DeleteCharacteristic(ctx context.Context, propertyID, characteristicID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteCharacteristic(ctx, propertyID, characteristicID)
}

// === REVIEWS ===

func (s *Service) ListReviews(ctx context.Context, page, perPage int) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListReviews(ctx, page, perPage)
}

func (s *Service) GetReview(ctx context.Context, reviewID string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetReview(ctx, reviewID)
}

func (s *Service) CreateReview(ctx context.Context, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateReview(ctx, args)
}

func (s *Service) UpdateReview(ctx context.Context, reviewID string, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.UpdateReview(ctx, reviewID, args)
}

func (s *Service) DeleteReview(ctx context.Context, reviewID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteReview(ctx, reviewID)
}

// === FILES ===

func (s *Service) ListFiles(ctx context.Context) ([]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.ListFiles(ctx)
}

func (s *Service) GetFile(ctx context.Context, fileID string) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetFile(ctx, fileID)
}

func (s *Service) CreateFile(ctx context.Context, args map[string]any) (any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return client.CreateFile(ctx, args)
}

func (s *Service) DeleteFile(ctx context.Context, fileID string) error {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return client.DeleteFile(ctx, fileID)
}

// Generic CRUD service helpers for round 3+ resources

type crudClient interface {
	ListX(ctx context.Context) ([]any, error)
}

func (s *Service) readOnlyCheck(ctx context.Context) error {
	return s.checkReadOnly(ctx)
}

// === DOMAINS ===
func (s *Service) ListDomains(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListDomains(ctx)
}
func (s *Service) GetDomain(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetDomain(ctx, id)
}
func (s *Service) CreateDomain(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateDomain(ctx, args)
}
func (s *Service) UpdateDomain(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateDomain(ctx, id, args)
}
func (s *Service) DeleteDomain(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteDomain(ctx, id)
}

// === CLIENT GROUPS ===
func (s *Service) ListClientGroups(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListClientGroups(ctx)
}
func (s *Service) GetClientGroup(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetClientGroup(ctx, id)
}
func (s *Service) CreateClientGroup(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateClientGroup(ctx, args)
}
func (s *Service) UpdateClientGroup(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateClientGroup(ctx, id, args)
}
func (s *Service) DeleteClientGroup(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteClientGroup(ctx, id)
}

// === CUSTOM STATUSES ===
func (s *Service) ListCustomStatuses(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListCustomStatuses(ctx)
}
func (s *Service) GetCustomStatus(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetCustomStatus(ctx, id)
}
func (s *Service) CreateCustomStatus(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateCustomStatus(ctx, args)
}
func (s *Service) UpdateCustomStatus(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateCustomStatus(ctx, id, args)
}
func (s *Service) DeleteCustomStatus(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteCustomStatus(ctx, id)
}

// === PRODUCT FIELDS ===
func (s *Service) ListProductFields(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListProductFields(ctx)
}
func (s *Service) GetProductField(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetProductField(ctx, id)
}
func (s *Service) CreateProductField(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateProductField(ctx, args)
}
func (s *Service) UpdateProductField(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateProductField(ctx, id, args)
}
func (s *Service) DeleteProductField(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteProductField(ctx, id)
}

// === VARIANT FIELDS ===
func (s *Service) ListVariantFields(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListVariantFields(ctx)
}
func (s *Service) GetVariantField(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetVariantField(ctx, id)
}
func (s *Service) CreateVariantField(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateVariantField(ctx, args)
}
func (s *Service) UpdateVariantField(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateVariantField(ctx, id, args)
}
func (s *Service) DeleteVariantField(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteVariantField(ctx, id)
}

// === COLLECTION FIELDS ===
func (s *Service) ListCollectionFields(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListCollectionFields(ctx)
}
func (s *Service) GetCollectionField(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetCollectionField(ctx, id)
}
func (s *Service) CreateCollectionField(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateCollectionField(ctx, args)
}
func (s *Service) UpdateCollectionField(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateCollectionField(ctx, id, args)
}
func (s *Service) DeleteCollectionField(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteCollectionField(ctx, id)
}

// === COLLECTION FILTERS ===
func (s *Service) ListCollectionFilters(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListCollectionFilters(ctx)
}
func (s *Service) GetCollectionFilter(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetCollectionFilter(ctx, id)
}
func (s *Service) CreateCollectionFilter(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateCollectionFilter(ctx, args)
}
func (s *Service) UpdateCollectionFilter(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateCollectionFilter(ctx, id, args)
}
func (s *Service) DeleteCollectionFilter(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteCollectionFilter(ctx, id)
}

// === JS TAGS ===
func (s *Service) ListJsTags(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListJsTags(ctx)
}
func (s *Service) GetJsTag(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetJsTag(ctx, id)
}
func (s *Service) CreateJsTag(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateJsTag(ctx, args)
}
func (s *Service) UpdateJsTag(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateJsTag(ctx, id, args)
}
func (s *Service) DeleteJsTag(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteJsTag(ctx, id)
}

// === EXTERNAL DISCOUNTS ===
func (s *Service) ListExternalDiscounts(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListExternalDiscounts(ctx)
}
func (s *Service) GetExternalDiscount(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetExternalDiscount(ctx, id)
}
func (s *Service) CreateExternalDiscount(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateExternalDiscount(ctx, args)
}
func (s *Service) UpdateExternalDiscount(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateExternalDiscount(ctx, id, args)
}
func (s *Service) DeleteExternalDiscount(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteExternalDiscount(ctx, id)
}

// === DISCOUNT COLLECTIONS GROUPS ===
func (s *Service) ListDiscountCollectionsGroups(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListDiscountCollectionsGroups(ctx)
}
func (s *Service) GetDiscountCollectionsGroup(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetDiscountCollectionsGroup(ctx, id)
}
func (s *Service) CreateDiscountCollectionsGroup(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateDiscountCollectionsGroup(ctx, args)
}
func (s *Service) UpdateDiscountCollectionsGroup(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateDiscountCollectionsGroup(ctx, id, args)
}
func (s *Service) DeleteDiscountCollectionsGroup(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteDiscountCollectionsGroup(ctx, id)
}

// === WAREHOUSES ===
func (s *Service) ListWarehouses(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListWarehouses(ctx)
}
func (s *Service) GetWarehouse(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetWarehouse(ctx, id)
}
func (s *Service) CreateWarehouse(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateWarehouse(ctx, args)
}
func (s *Service) UpdateWarehouse(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateWarehouse(ctx, id, args)
}
func (s *Service) DeleteWarehouse(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteWarehouse(ctx, id)
}

// === SALES CHANNELS ===
func (s *Service) ListSalesChannels(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListSalesChannels(ctx)
}
func (s *Service) GetSalesChannel(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetSalesChannel(ctx, id)
}
func (s *Service) CreateSalesChannel(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateSalesChannel(ctx, args)
}
func (s *Service) UpdateSalesChannel(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateSalesChannel(ctx, id, args)
}
func (s *Service) DeleteSalesChannel(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteSalesChannel(ctx, id)
}

// === STOCK CURRENCIES ===
func (s *Service) ListStockCurrencies(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListStockCurrencies(ctx)
}
func (s *Service) GetStockCurrency(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetStockCurrency(ctx, id)
}
func (s *Service) CreateStockCurrency(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateStockCurrency(ctx, args)
}
func (s *Service) UpdateStockCurrency(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateStockCurrency(ctx, id, args)
}
func (s *Service) DeleteStockCurrency(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteStockCurrency(ctx, id)
}

// === PRICE KINDS ===
func (s *Service) ListPriceKinds(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListPriceKinds(ctx)
}
func (s *Service) GetPriceKind(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetPriceKind(ctx, id)
}
func (s *Service) CreatePriceKind(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreatePriceKind(ctx, args)
}
func (s *Service) UpdatePriceKind(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdatePriceKind(ctx, id, args)
}
func (s *Service) DeletePriceKind(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeletePriceKind(ctx, id)
}

// === OPTION NAMES ===
func (s *Service) ListOptionNames(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListOptionNames(ctx)
}
func (s *Service) GetOptionName(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetOptionName(ctx, id)
}
func (s *Service) CreateOptionName(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateOptionName(ctx, args)
}
func (s *Service) UpdateOptionName(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateOptionName(ctx, id, args)
}
func (s *Service) DeleteOptionName(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteOptionName(ctx, id)
}

// === OPTION VALUES ===
func (s *Service) ListOptionValues(ctx context.Context, optionNameID string) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListOptionValues(ctx, optionNameID)
}
func (s *Service) GetOptionValue(ctx context.Context, optionNameID, optionValueID string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetOptionValue(ctx, optionNameID, optionValueID)
}
func (s *Service) CreateOptionValue(ctx context.Context, optionNameID string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateOptionValue(ctx, optionNameID, args)
}
func (s *Service) UpdateOptionValue(ctx context.Context, optionNameID, optionValueID string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateOptionValue(ctx, optionNameID, optionValueID, args)
}
func (s *Service) DeleteOptionValue(ctx context.Context, optionNameID, optionValueID string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteOptionValue(ctx, optionNameID, optionValueID)
}

// === APPLICATION ACTIONS ===
func (s *Service) ListApplicationActions(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListApplicationActions(ctx)
}
func (s *Service) GetApplicationAction(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetApplicationAction(ctx, id)
}
func (s *Service) CreateApplicationAction(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateApplicationAction(ctx, args)
}
func (s *Service) UpdateApplicationAction(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateApplicationAction(ctx, id, args)
}
func (s *Service) DeleteApplicationAction(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteApplicationAction(ctx, id)
}

// === APPLICATION CHARGES ===
func (s *Service) ListApplicationCharges(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListApplicationCharges(ctx)
}
func (s *Service) GetApplicationCharge(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetApplicationCharge(ctx, id)
}
func (s *Service) CreateApplicationCharge(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateApplicationCharge(ctx, args)
}
func (s *Service) DeclineApplicationCharge(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeclineApplicationCharge(ctx, id)
}

// === APPLICATION WIDGETS ===
func (s *Service) ListApplicationWidgets(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListApplicationWidgets(ctx)
}
func (s *Service) GetApplicationWidget(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetApplicationWidget(ctx, id)
}
func (s *Service) CreateApplicationWidget(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateApplicationWidget(ctx, args)
}
func (s *Service) UpdateApplicationWidget(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateApplicationWidget(ctx, id, args)
}
func (s *Service) DeleteApplicationWidget(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteApplicationWidget(ctx, id)
}

// === RECURRING APPLICATION CHARGES ===
func (s *Service) ListRecurringApplicationCharges(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListRecurringApplicationCharges(ctx)
}
func (s *Service) GetRecurringApplicationCharge(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetRecurringApplicationCharge(ctx, id)
}
func (s *Service) CreateRecurringApplicationCharge(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateRecurringApplicationCharge(ctx, args)
}
func (s *Service) DeleteRecurringApplicationCharge(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteRecurringApplicationCharge(ctx, id)
}

// === BONUS SYSTEM TRANSACTIONS ===
func (s *Service) ListBonusSystemTransactions(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListBonusSystemTransactions(ctx)
}
func (s *Service) GetBonusSystemTransaction(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetBonusSystemTransaction(ctx, id)
}
func (s *Service) CreateBonusSystemTransaction(ctx context.Context, clientID int, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateBonusSystemTransaction(ctx, clientID, args)
}
func (s *Service) ListClientBonusSystemTransactions(ctx context.Context, clientID string) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListClientBonusSystemTransactions(ctx, clientID)
}

// === OAUTH APPS ===
func (s *Service) ListOAuthApps(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListOAuthApps(ctx)
}
func (s *Service) GetOAuthApp(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetOAuthApp(ctx, id)
}
func (s *Service) CreateOAuthApp(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateOAuthApp(ctx, args)
}
func (s *Service) UpdateOAuthApp(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateOAuthApp(ctx, id, args)
}
func (s *Service) DeleteOAuthApp(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteOAuthApp(ctx, id)
}

// === USER TASKS ===
func (s *Service) ListUserTasks(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListUserTasks(ctx)
}
func (s *Service) GetUserTask(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetUserTask(ctx, id)
}
func (s *Service) CreateUserTask(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateUserTask(ctx, args)
}
func (s *Service) UpdateUserTask(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateUserTask(ctx, id, args)
}
func (s *Service) DeleteUserTask(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteUserTask(ctx, id)
}

// === SUPPLEMENTARY ===
func (s *Service) ListSupplementaries(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListSupplementaries(ctx)
}
func (s *Service) GetSupplementary(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetSupplementary(ctx, id)
}
func (s *Service) CreateSupplementary(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateSupplementary(ctx, args)
}
func (s *Service) DeleteSupplementary(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteSupplementary(ctx, id)
}

// === SIMILARS ===
func (s *Service) ListSimilars(ctx context.Context, productID string) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListSimilars(ctx, productID)
}
func (s *Service) CreateSimilar(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateSimilar(ctx, args)
}
func (s *Service) DeleteSimilar(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteSimilar(ctx, id)
}

// === UNITS ===
func (s *Service) ListUnits(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListUnits(ctx)
}

// === PIPELINE TOOLS ===

func (s *Service) SearchProductDetails(ctx context.Context, query string, productID string) (map[string]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	var product *insales.Product
	if productID != "" {
		product, err = client.GetProduct(ctx, productID)
	} else {
		products, err2 := client.ListProducts(ctx, 1, 5, map[string]any{"q": query})
		if err2 != nil {
			return nil, err2
		}
		if len(products) == 0 {
			return map[string]any{"products": []any{}, "message": "Товары не найдены"}, nil
		}
		product = &products[0]
	}
	if err != nil {
		return nil, err
	}

	// Get meta fields for this product
	metaFields, _ := client.ListMetaFields(ctx, map[string]any{
		"resource_type": "Product",
		"resource_id":   product.ID,
	})

	result := map[string]any{
		"product":     product,
		"meta_fields": metaFields,
	}
	return result, nil
}

func (s *Service) OrderFullDetails(ctx context.Context, orderID string) (map[string]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	order, err := client.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	var clientData any
	if order.ClientID > 0 {
		clientData, _ = client.GetClient(ctx, fmt.Sprintf("%d", order.ClientID))
	}

	// Get meta fields for this order
	metaFields, _ := client.ListMetaFields(ctx, map[string]any{
		"resource_type": "Order",
		"resource_id":   order.ID,
	})

	result := map[string]any{
		"order":       order,
		"client":      clientData,
		"meta_fields": metaFields,
	}
	return result, nil
}

func (s *Service) ClientFullProfile(ctx context.Context, clientID string) (map[string]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	clientData, err := client.GetClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Get orders for this client
	orders, _ := client.ListOrders(ctx, 1, 10, map[string]any{})

	// Get meta fields for this client
	metaFields, _ := client.ListMetaFields(ctx, map[string]any{
		"resource_type": "Client",
		"resource_id":   clientID,
	})

	// Get bonus transactions
	bonusTxns, _ := client.ListClientBonusSystemTransactions(ctx, clientID)

	result := map[string]any{
		"client":             clientData,
		"recent_orders":      orders,
		"meta_fields":        metaFields,
		"bonus_transactions": bonusTxns,
	}
	return result, nil
}

func (s *Service) ProductCatalogSearch(ctx context.Context, query string, filters map[string]any) (map[string]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	page := 1
	perPage := 25
	if v, ok := filters["page"].(float64); ok && v > 0 {
		page = int(v)
	}
	if v, ok := filters["per_page"].(float64); ok && v > 0 {
		perPage = int(v)
	}

	searchFilters := map[string]any{
		"page":     page,
		"per_page": perPage,
	}
	if query != "" {
		searchFilters["q"] = query
	}
	if v, ok := filters["category_id"].(string); ok && v != "" {
		searchFilters["category_id"] = v
	}
	if v, ok := filters["collection_id"].(string); ok && v != "" {
		searchFilters["collection_id"] = v
	}

	products, err := client.ListProducts(ctx, page, perPage, searchFilters)
	if err != nil {
		return nil, err
	}

	result := map[string]any{
		"products": products,
		"page":     page,
		"per_page": perPage,
	}
	return result, nil
}

// === OUTLETS ===
func (s *Service) ListOutlets(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListOutlets(ctx)
}
func (s *Service) GetOutlet(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetOutlet(ctx, id)
}
func (s *Service) CreateOutlet(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateOutlet(ctx, args)
}
func (s *Service) UpdateOutlet(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateOutlet(ctx, id, args)
}
func (s *Service) DeleteOutlet(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteOutlet(ctx, id)
}

// === MARKETPLACES ===
func (s *Service) ListMarketplaces(ctx context.Context) ([]any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListMarketplaces(ctx)
}
func (s *Service) GetMarketplace(ctx context.Context, id string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetMarketplace(ctx, id)
}
func (s *Service) CreateMarketplace(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateMarketplace(ctx, args)
}
func (s *Service) UpdateMarketplace(ctx context.Context, id string, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateMarketplace(ctx, id, args)
}
func (s *Service) DeleteMarketplace(ctx context.Context, id string) error {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return err
	}
	return c.DeleteMarketplace(ctx, id)
}

// === PROVIDER TOKEN ===
func (s *Service) CreateProviderToken(ctx context.Context, messengerType string) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.CreateProviderToken(ctx, messengerType)
}

// === ACCOUNT ===
func (s *Service) UpdateAccount(ctx context.Context, args map[string]any) (any, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.checkReadOnly(ctx); err != nil {
		return nil, err
	}
	return c.UpdateAccount(ctx, args)
}

// === COUNTS ===
func (s *Service) GetOrdersCount(ctx context.Context, filters map[string]any) (int, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return 0, err
	}
	return c.GetOrdersCount(ctx, filters)
}
func (s *Service) GetProductsCount(ctx context.Context) (int, error) {
	c, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return 0, err
	}
	return c.GetProductsCount(ctx)
}

// === PIPELINE: CollectionProducts ===
func (s *Service) CollectionProducts(ctx context.Context, collectionID string) (map[string]any, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	collection, err := client.GetCollection(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("get collection: %w", err)
	}

	products, err := client.ListProducts(ctx, 1, 250, map[string]any{"collection_id": collectionID})
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	return map[string]any{
		"collection": collection,
		"products":   products,
	}, nil
}

// === DATA EXPORT (CSV) ===

func (s *Service) ExportProducts(ctx context.Context) (string, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return "", err
	}

	headers := []string{"id", "title", "price", "sku", "vendor", "quantity", "created_at"}
	var rows [][]string

	page := 1
	for {
		products, err := client.ListProducts(ctx, page, 250, nil)
		if err != nil {
			return "", fmt.Errorf("list products page %d: %w", page, err)
		}
		if len(products) == 0 {
			break
		}
		for _, p := range products {
			rows = append(rows, []string{
				p.ID, p.Title, p.Price, p.SKU, p.Vendor,
				fmt.Sprintf("%d", p.Quantity), p.CreatedAt,
			})
		}
		if len(products) < 250 {
			break
		}
		page++
	}

	return formatCSV(headers, rows), nil
}

func (s *Service) ExportOrders(ctx context.Context) (string, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return "", err
	}

	headers := []string{"id", "number", "fulfillment_status", "total_price", "client_name", "created_at"}
	var rows [][]string

	page := 1
	for {
		orders, err := client.ListOrders(ctx, page, 250, nil)
		if err != nil {
			return "", fmt.Errorf("list orders page %d: %w", page, err)
		}
		if len(orders) == 0 {
			break
		}
		for _, o := range orders {
			rows = append(rows, []string{
				o.ID, o.Number, o.FulfillmentStatus, o.TotalPrice,
				o.ClientName, o.CreatedAt,
			})
		}
		if len(orders) < 250 {
			break
		}
		page++
	}

	return formatCSV(headers, rows), nil
}

func (s *Service) ExportClients(ctx context.Context) (string, error) {
	client, ctx, err := s.manager.GetClient(ctx)
	if err != nil {
		return "", err
	}

	headers := []string{"id", "name", "surname", "email", "phone", "created_at"}
	var rows [][]string

	page := 1
	for {
		clients, err := client.ListClients(ctx, page, 250)
		if err != nil {
			return "", fmt.Errorf("list clients page %d: %w", page, err)
		}
		if len(clients) == 0 {
			break
		}
		for _, c := range clients {
			m, ok := c.(map[string]any)
			if !ok {
				continue
			}
			rows = append(rows, []string{
				fmt.Sprintf("%v", m["id"]),
				fmt.Sprintf("%v", m["name"]),
				fmt.Sprintf("%v", m["surname"]),
				fmt.Sprintf("%v", m["email"]),
				fmt.Sprintf("%v", m["phone"]),
				fmt.Sprintf("%v", m["created_at"]),
			})
		}
		if len(clients) < 250 {
			break
		}
		page++
	}

	return formatCSV(headers, rows), nil
}

func formatCSV(headers []string, rows [][]string) string {
	result := ""
	for i, h := range headers {
		if i > 0 {
			result += ","
		}
		result += csvEscape(h)
	}
	result += "\n"
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				result += ","
			}
			result += csvEscape(cell)
		}
		result += "\n"
	}
	return result
}

func csvEscape(s string) string {
	if s == "" || s == "<nil>" {
		return ""
	}
	needsQuote := false
	for _, c := range s {
		if c == ',' || c == '"' || c == '\n' || c == '\r' {
			needsQuote = true
			break
		}
	}
	if needsQuote {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}
